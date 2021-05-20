package gorilla_mux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/silverspase/k8s-prod-service/internal/todo"
	"github.com/silverspase/k8s-prod-service/internal/todo/model"
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
	t.logger.Info("transport.CreateItems")
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
	t.logger.Info("transport.GetAllItems")
	ctx := context.Background()

	items, err := t.useCase.GetAllItems(ctx)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

func (t *transport) GetItem(w http.ResponseWriter, r *http.Request) {
	t.logger.Info("transport.GetItem")
	ctx := context.Background()

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missed id path param"})
		return
	}
	item, ok := t.useCase.GetItem(ctx, id)
	if !ok {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("item with id %v not found", id)})
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

func (t *transport) UpdateItem(w http.ResponseWriter, r *http.Request) {
	t.logger.Info("transport.UpdateItem")
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
	t.logger.Info("transport.DeleteItem")
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
