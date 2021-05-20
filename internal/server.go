package internal

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	meta "github.com/silverspase/k8s-prod-service/internal/metadata/transport/gorilla-mux"
	"github.com/silverspase/k8s-prod-service/internal/todo"
)

type server struct {
	transport todo.Transport
	logger    *zap.Logger
	isReady   *atomic.Value
}

func NewServer(logger *zap.Logger, transport todo.Transport) *server {
	return &server{
		transport: transport,
		logger:    logger,
		isReady:   &atomic.Value{},
	}
}

// Router registers routes and returns an instance of a router.
func (s server) GorillaMuxRouter(buildTime, commit, release string) *mux.Router {
	go func() {
		log.Printf("wait for Ready probe...")
		time.Sleep(2 * time.Second)
		s.isReady.Store(true)
		log.Printf("Ready probe is positive.")
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", meta.GetMetaData(buildTime, commit, release)).Methods(http.MethodGet)
	r.HandleFunc("/health", meta.HealthCheck)
	r.HandleFunc("/readiness", meta.Readiness(s.isReady))

	todoR := r.PathPrefix("/todo").Subrouter()
	todoR.Path("/").HandlerFunc(s.transport.CreateItem).Methods(http.MethodPost)
	todoR.Path("/").HandlerFunc(s.transport.GetAllItems).Methods(http.MethodGet)
	todoR.Path("/{id}").HandlerFunc(s.transport.GetItem).Methods(http.MethodGet)
	todoR.Path("/{id}").HandlerFunc(s.transport.UpdateItem).Methods(http.MethodPut)
	todoR.Path("/{id}").HandlerFunc(s.transport.DeleteItem).Methods(http.MethodDelete)

	return r
}
