package gorilla_mux

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

type metaData struct {
	BuildTime string `json:"buildTime"`
	Commit    string `json:"commit"`
	Release   string `json:"release"`
}

// healthCheck is a live-ness probe.
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readiness probe.
func Readiness(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// GetMetaData returns service info.
func GetMetaData(buildTime, commit, release string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		info := metaData{buildTime, commit, release}

		body, err := json.Marshal(info)
		if err != nil {
			log.Printf("Could not encode info data: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
