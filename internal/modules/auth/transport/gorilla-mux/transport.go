package gorilla_mux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/silverspase/todo/internal/modules/auth"
	"github.com/silverspase/todo/internal/modules/auth/model"
)

type transport struct {
	useCase auth.UseCase
	logger  *zap.Logger
}

func NewTransport(logger *zap.Logger, useCase auth.UseCase) auth.Transport {
	return &transport{
		useCase: useCase,
		logger:  logger,
	}
}

func (t *transport) CreateUser(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("CreateUser")
	ctx := context.Background()
	defer r.Body.Close()

	var user model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	id, err := t.useCase.CreateUser(ctx, user)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "created", "id": id})
}

func (t *transport) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("GetAllUsers")
	ctx := context.Background()

	var page int
	var err error

	pageStr := r.FormValue("page")
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "page param is not a number"})
			return
		}
	}

	users, err := t.useCase.GetAllUsers(ctx, page)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (t *transport) GetUser(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("GetUser")
	ctx := context.Background()

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missed id path param"})
		return
	}

	item, err := t.useCase.GetUser(ctx, id)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get entry with id %v: %v", id, err)})
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

func (t *transport) UpdateUser(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("UpdateUser")
	ctx := context.Background()
	defer r.Body.Close()

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missed id path param"})
		return
	}

	var item model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"unmarshal error": "Invalid request payload"})
		return
	}

	item.ID = id
	id, err := t.useCase.UpdateUser(ctx, item)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "updated", "id": id})
}

func (t *transport) DeleteUser(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("DeleteUser")
	ctx := context.Background()

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missed id path param"})
		return
	}
	id, err := t.useCase.DeleteUser(ctx, id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("item with id %v not found", id)})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
