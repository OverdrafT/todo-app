package gorilla_mux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/silverspase/todo/internal/todo"
	"github.com/silverspase/todo/internal/todo/model"
)

type transport struct {
	useCase todo.UseCase
	logger  *zap.Logger
}

func NewTransport(logger *zap.Logger, useCase todo.UseCase) todo.Transport {
	return &transport{
		useCase: useCase,
		logger:  logger,
	}
}

func (t *transport) CreateItem(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("transport.CreateItem")
	ctx := context.Background()
	defer r.Body.Close()

	var item model.Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	id, err := t.useCase.CreateItem(ctx, item)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "created", "id": id})
}

func (t *transport) GetAllItems(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("GetAllItems")
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

	items, err := t.useCase.GetAllItems(ctx, page)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

func (t *transport) GetItem(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("GetItem")
	ctx := context.Background()

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missed id path param"})
		return
	}
	item, err := t.useCase.GetItem(ctx, id)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("item with id %v not found", id)})
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

func (t *transport) UpdateItem(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("UpdateItem")
	ctx := context.Background()
	defer r.Body.Close()

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missed id path param"})
		return
	}

	var item model.Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"unmarshal error": "Invalid request payload"})
		return
	}

	item.ID = id
	id, err := t.useCase.UpdateItem(ctx, item)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "updated", "id": id})
}

func (t *transport) DeleteItem(w http.ResponseWriter, r *http.Request) {
	t.logger.Debug("DeleteItem")
	ctx := context.Background()

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missed id path param"})
		return
	}
	id, err := t.useCase.DeleteItem(ctx, id)
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
