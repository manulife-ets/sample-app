package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"git.platform.manulife.io/oss/url-shortener/db"
	"git.platform.manulife.io/oss/url-shortener/routes"
	"git.platform.manulife.io/oss/url-shortener/utils"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
)

const (
	// DefaultLocalPort the default port when OS ENV is not supplied. Typically only happens on local
	DefaultLocalPort = "9000"
)

// BASIC AUTH
var username, password string
var secured = false

func init() {
	env := os.Getenv("LOCAL")
	if len(env) > 1 {
		username = "admin"
		password = "god"
	} else {
		username = os.Getenv("BASIC_AUTH_USER")
		if len(username) < 1 {
			fmt.Println("Missing the Basic Auth Username in environment")
			os.Exit(0)
		}
		password = os.Getenv("BASIC_AUTH_PASS")
		if len(password) < 1 {
			fmt.Println("Missing the Basic Auth Password in environment")
			os.Exit(0)
		}
	}

	// Undeclared ENV variable will result in NO Security enabled
	secured, _ = strconv.ParseBool(os.Getenv("SECURED"))
}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		fmt.Printf("Warning, PORT not set. Defaulting to %v\n", DefaultLocalPort)
		port = DefaultLocalPort
	}

	err := db.InitDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	router := mux.NewRouter()

	router.HandleFunc(newrelic.WrapHandleFunc(utils.NR(), "/api", routes.SwaggerEndpoint)).Methods("GET")

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// /v2
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// INSERT
	router.HandleFunc(newrelic.WrapHandleFunc(utils.NR(), "/api/v2", basicAuth(routes.InsertEndpoint, "Please enter your username and password for this site"))).Methods("POST")
	// UPDATE
	router.HandleFunc(newrelic.WrapHandleFunc(utils.NR(), "/api/v2/{id}", basicAuth(routes.UpdateEndpoint, "Please enter your username and password for this site"))).Methods("PUT")

	// Find All
	router.HandleFunc(newrelic.WrapHandleFunc(utils.NR(), "/api/v2", routes.FindAllEndpoint)).Methods("GET")
	// Find One
	router.HandleFunc(newrelic.WrapHandleFunc(utils.NR(), "/api/v2/{id}", routes.FindOneEndpoint)).Methods("GET")

	// FORWARD
	router.HandleFunc(newrelic.WrapHandleFunc(utils.NR(), "/{id}", routes.RouterEndpoint)).Methods("GET")

	// DELETE
	router.HandleFunc(newrelic.WrapHandleFunc(utils.NR(), "/api/v2", basicAuth(routes.DeleteAllEndpoint, "Please enter your username and password for this site"))).Methods("DELETE")

	// Enabled only for performance testing
	// router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	// Serve index.html static page
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	// fmt.Println(http.ListenAndServe(":"+port, router))
	var wait = time.Second * 15
	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("Gracefully shutting down")
	os.Exit(0)
}

// Temporarily protect end points with basic auth, eventually move to Client_Credentials
func basicAuth(handler http.HandlerFunc, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if secured {
			user, pass, ok := r.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorized.\n"))
				return
			}
		}
		handler(w, r)
	}
}
