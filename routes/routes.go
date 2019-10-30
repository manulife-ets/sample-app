package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"git.platform.manulife.io/oss/url-shortener/db"
	"git.platform.manulife.io/oss/url-shortener/model"
	"git.platform.manulife.io/oss/url-shortener/utils"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	newrelic "github.com/newrelic/go-agent"
)

// SwaggerEndpoint get to longurl based on requested id
func SwaggerEndpoint(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	defer req.Body.Close()

	f, err := os.Open("./SWAGGER/swagger.html")
	if err != nil {
		fmt.Printf("file error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer f.Close() // ensure the file is closed

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if _, err := io.Copy(w, f); err != nil {
		fmt.Printf("file copy error: %v\n", err)
	}
}

// RouterEndpoint get to longurl based on requested id
func RouterEndpoint(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	defer req.Body.Close()

	url, err := validateIncomingID(w, req)
	if !err || (model.URL{}) == *url {
		return
	}

	redirect := req.URL.Query().Get("redirect")
	if len(redirect) > 0 {
		// Only handle false, else redirect!
		if redirect == "false" {
			json.NewEncoder(w).Encode(url)
			return
		}
	}

	// Track Metrics
	go func() {
		db.Upsert(txn,
			"urlmappings",
			bson.D{{"_id", url.ID}},
			bson.D{{"$inc", bson.D{{"metrics.routedcount", 1}}}})
	}()

	// Default to redirect
	http.Redirect(w, req, url.LongURL, http.StatusTemporaryRedirect)
}

// FindOneEndpoint get to longurl based on requested id
func FindOneEndpoint(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	defer req.Body.Close()

	url, err := validateIncomingID(w, req)
	if !err || (model.URL{}) == *url {
		return
	}

	redirect := req.URL.Query().Get("redirect")
	if len(redirect) > 0 {
		// Only handle true to redirect!
		if redirect == "true" {
			http.Redirect(w, req, url.LongURL, http.StatusTemporaryRedirect)
			return
		}
	}
	// Default to redirect
	json.NewEncoder(w).Encode(url)
}

// FindAllEndpoint find all, no real usecase for this yet!
func FindAllEndpoint(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	defer req.Body.Close()

	urls, err := db.Find(txn, "urlmappings", bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(w).Encode(urls)
}

// InsertEndpoint POST:
// Shorten a long url, optionally providing a custom id.
// {
//	  id: optional
//    longurl: mandatory
// }
//
// 201: Insert successful
// 400: Missing longurl or invalid format longurl
// 409: Duplicate ID issue
// 500: Mongo DB Collection issues
func InsertEndpoint(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	defer req.Body.Close()

	var url model.URL
	_ = json.NewDecoder(req.Body).Decode(&url)

	if len(url.LongURL) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Missing URL to Shorten" }`))
		return
	}

	if !utils.ValidateURLFormat([]byte(url.LongURL)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Invalid format for URL to Shorten. Pattern> http[s]://[www.]google.com" }`))
		return
	}
	if !utils.ValidateURLXSS([]byte(url.LongURL)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Invalid content for URL to Shorten. May contain malicious content." }`))
		return
	}

	if len(url.ID) < 1 {
		id, err := utils.GenerateID()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		url.ID = id
	}

	// Check if ID already exists, post will be rejected
	count, err := db.Count(txn, "urlmappings", bson.M{"_id": url.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{ "message": Duplicate ID [` + url.ID + `], insert rejected" }`))
		return
	}

	ts := time.Now().UTC()
	url.CreatedAt = &ts
	url.CreatedBy = "SYSTEM"

	resp := db.InsertOne(txn, "urlmappings", url)
	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "Unable to store record" }`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(url)
}

// UpdateEndpoint PUT:
// Update a long url by supplying the {id}.
// {
//    longurl: mandatory
// }
//
// 202: Update successful
// 400: Missing longurl or Invalid format longurl
// 404: ID to update NOT found
// 500: Mongo DB Collection issues
func UpdateEndpoint(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	defer req.Body.Close()

	var url model.URL
	_ = json.NewDecoder(req.Body).Decode(&url)

	if len(url.LongURL) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Missing URL to Shorten" }`))
		return
	}

	if !utils.ValidateURLFormat([]byte(url.LongURL)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Invalid format for URL to Shorten. Pattern> http[s]://[www.]google.com" }`))
		return
	}
	if !utils.ValidateURLXSS([]byte(url.LongURL)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Invalid content for URL to Shorten. May contain malicious content." }`))
		return
	}

	params := mux.Vars(req)
	id := params["id"]
	url.ID = id

	// Check if ID already exists, put will be rejected
	count, err := db.Count(txn, "urlmappings", bson.M{"_id": url.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if count < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{ "message": ID [` + url.ID + `] not found, update rejected" }`))
		return
	}

	ts := time.Now().UTC()
	url.UpdatedAt = &ts
	url.UpdatedBy = "SYSTEM"

	filter := bson.D{{"_id", url.ID}}

	update := bson.D{{"$set", bson.D{{"longurl", url.LongURL}, {"updatedat", url.UpdatedAt}, {"updatedby", url.UpdatedBy}}}}

	_, err = db.Upsert(txn, "urlmappings", filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(url)
}

// DeleteAllEndpoint delete all, no real usecase for this yet!
func DeleteAllEndpoint(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	defer req.Body.Close()

	count, err := db.Delete(txn, "urlmappings", bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "Deleted [` + strconv.FormatInt(count, 10) + `], documents" }`))
	return
}

// validateIncomingID return a FALSE if an issue occurred validating incoming
func validateIncomingID(w http.ResponseWriter, req *http.Request) (*model.URL, bool) {
	txn := newrelic.FromContext(req.Context())
	defer txn.End()

	params := mux.Vars(req)
	id := params["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Missing ID" }`))
		return &model.URL{}, false
	}

	url := &model.URL{}

	db.FindOne(txn, "urlmappings", bson.M{"_id": id}, url)
	if (model.URL{}) == *url {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{ "message": "Invalid ID" }`))
		return &model.URL{}, false
	}
	return url, true
}
