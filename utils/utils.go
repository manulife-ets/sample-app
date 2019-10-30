package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"git.platform.manulife.io/go-common/pcf/monitoring"
	newrelic "github.com/newrelic/go-agent"
	"github.com/speps/go-hashids"
)

var nr newrelic.Application

// https://regex101.com/r/Xv3sto/1
var urlPattern = `^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`
var rURL = regexp.MustCompile(urlPattern)

// https://regex101.com/r/9wGPP1/2
var urlXSSPattern = `(?mi)(\b)(on\S+)(\s*)=|javascript|(<\s*)(\/*)script`
var rURLXSS = regexp.MustCompile(urlXSSPattern)

func init() {
	var err error

	env := os.Getenv("LOCAL")
	if len(env) > 1 {
		cfg := newrelic.NewConfig("url-shortener-local", os.Getenv("NR_LICENSE"))
		cfg.CustomInsightsEvents.Enabled = true
		// Disable NR on local
		// cfg.Enabled = false
		nr, err = newrelic.NewApplication(cfg)
	} else {
		nr, err = monitoring.InitNewRelic()
	}
	if err != nil {
		fmt.Println(err.Error())
	}

}

// NR return newrelic app
func NR() newrelic.Application {
	return nr
}

// TimeTrack track elapsed time
func TimeTrack(start time.Time, name string, code string) {
	elapsed := time.Since(start)
	nr.RecordCustomMetric(code, elapsed.Seconds()/1000)
	fmt.Printf("%s took %s\n", name, elapsed)
}

// GenerateID generate a unique id
func GenerateID() (string, error) {
	hd := hashids.NewData()
	hd.Alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	hashedID, err := h.Encode([]int{int(now.UnixNano())})
	if err != nil {
		fmt.Printf("Error hash encoding %v\n", err.Error())
		return "", errors.New("Unable to create uniqe hash")
	}
	return hashedID, nil
}

// ValidateURLFormat validate the format
func ValidateURLFormat(needle []byte) (ok bool) {
	return len(rURL.FindAllSubmatch(needle, -1)) > 0
}

// ValidateURLXSS validate the format
func ValidateURLXSS(needle []byte) (ok bool) {
	return len(rURLXSS.FindAllSubmatch(needle, -1)) < 1
}
