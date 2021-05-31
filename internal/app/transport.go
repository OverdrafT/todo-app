package app

import (
	"net/http"

	"github.com/gorilla/mux"

	meta "github.com/silverspase/todo/internal/modules/metadata/transport/gorilla-mux"
)

// TODO move router init to separate package (resolve cycle import issue)
func gorillaMuxRouter(t *App) http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/health", meta.HealthCheck)
	// r.HandleFunc("/readiness", meta.Readiness(s.isReady))

	todo := r.PathPrefix("/todo").Subrouter()
	todo.Path("/").HandlerFunc(t.Todo.CreateItem).Methods(http.MethodPost)
	todo.Path("/").HandlerFunc(t.Todo.GetAllItems).Methods(http.MethodGet)
	todo.Path("/{id}").HandlerFunc(t.Todo.GetItem).Methods(http.MethodGet)
	todo.Path("/{id}").HandlerFunc(t.Todo.UpdateItem).Methods(http.MethodPut)
	todo.Path("/{id}").HandlerFunc(t.Todo.DeleteItem).Methods(http.MethodDelete)

	user := r.PathPrefix("/user").Subrouter()
	user.Path("/").HandlerFunc(t.Auth.CreateUser).Methods(http.MethodPost)
	user.Path("/").HandlerFunc(t.Auth.GetAllUsers).Methods(http.MethodGet)
	user.Path("/{id}").HandlerFunc(t.Auth.GetUser).Methods(http.MethodGet)
	user.Path("/{id}").HandlerFunc(t.Auth.UpdateUser).Methods(http.MethodPut)
	user.Path("/{id}").HandlerFunc(t.Auth.DeleteUser).Methods(http.MethodDelete)

	return r
}
