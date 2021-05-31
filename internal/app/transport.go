package app

import (
	"net/http"

	"github.com/gorilla/mux"

	meta "github.com/silverspase/todo/internal/metadata/transport/gorilla-mux"
)

// TODO move router init to separate package (resolve cycle import issue)
func gorillaMuxRouter(t *App) http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/health", meta.HealthCheck)
	// r.HandleFunc("/readiness", meta.Readiness(s.isReady))

	todoR := r.PathPrefix("/todo").Subrouter()
	todoR.Path("/").HandlerFunc(t.TodoT.CreateItem).Methods(http.MethodPost)
	todoR.Path("/").HandlerFunc(t.TodoT.GetAllItems).Methods(http.MethodGet)
	todoR.Path("/{id}").HandlerFunc(t.TodoT.GetItem).Methods(http.MethodGet)
	todoR.Path("/{id}").HandlerFunc(t.TodoT.UpdateItem).Methods(http.MethodPut)
	todoR.Path("/{id}").HandlerFunc(t.TodoT.DeleteItem).Methods(http.MethodDelete)

	return r
}
