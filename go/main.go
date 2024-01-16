package main

import (
	"crypto/tls"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	dayswarning := DaysWarning()
	timetopanic := false
	hostnames := GetExternalHostnames()
	sort.Strings(hostnames)
	log.Infof("cert-scanner scanning for certs within %v days of expiration...", dayswarning)
	for _, hostname := range hostnames {
		expirationdate, err := GetExpirationDate(hostname)
		if err != nil {
			log.Errorf("cert-scanner unable to process %v: %v", hostname, err)
			continue
		}
		daysleft := DaysLeft(expirationdate)
		if daysleft < dayswarning {
			timetopanic = true
			log.Errorf("cert-scanner failed: %v expires in %v days", hostname, daysleft)
		} else {
			log.Infof("cert-scanner passed: %v (%v days left)", hostname, daysleft)
		}
	}
	if timetopanic {
		// exit with a nonzero code if any impending expirations were found
		log.Errorf("cert-scanner failed: one or more certs require attention")
		os.Exit(1)
	}
	// otherwise signal clean completion in logs and exit cleanly
	log.Infof("cert-scanner passed")
}

// DaysWarning returns the number of days out a renewal can be before triggering an error.
func DaysWarning() int64 {
	d, err := strconv.ParseInt(os.Getenv("DAYS_WARNING"), 10, 64)
	if err != nil {
		panic(err)
	}
	return d
}

// GetExpirationDate returns an SSL certificate expiration date given a hostname.
func GetExpirationDate(h string) (time.Time, error) {
	conn, err := tls.Dial("tcp", h+":443", nil)
	if err != nil {
		return time.Now(), err
	}
	err = conn.VerifyHostname(h)
	if err != nil {
		return time.Now(), err
	}
	return conn.ConnectionState().PeerCertificates[0].NotAfter, nil
}

// DaysLeft returns an integer number of days until a given time.
func DaysLeft(t1 time.Time) int64 {
	return int64(math.Round(t1.Sub(time.Now()).Hours() / 24))
}

// GetExternalHostnames returns a list of every hostname present in a static hosts.txt file checked into the repo root.
func GetExternalHostnames() []string {
	content, err := ioutil.ReadFile("hosts.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(content), "\n")
}
