package handlers

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// Router registers routes and returns an instance of a router.
func Router(buildTime, commit, release string) *mux.Router {
	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Printf("wait for Ready probe...")
		time.Sleep(5 * time.Second)
		isReady.Store(true)
		log.Printf("Ready probe is positive.")
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", getMetaData(buildTime, commit, release)).Methods("GET")
	r.HandleFunc("/health", healthCheck)
	r.HandleFunc("/readiness", readiness(isReady))

	return r
}
