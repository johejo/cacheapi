package cacheapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/allegro/bigcache"
	"github.com/go-chi/chi"
	"golang.org/x/xerrors"
)

// CacheHandler combines several methods those implement http.Handler.
type CacheHandler interface {
	PutItem(w http.ResponseWriter, r *http.Request)
	PutItemWithKey(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	ListItem(w http.ResponseWriter, r *http.Request)
	GetStats(w http.ResponseWriter, r *http.Request)
	GetStatus(w http.ResponseWriter, r *http.Request)
}

type cacheHandler struct {
	CacheService
}

// NewCacheHandler returns a new CacheHandler.
func NewCacheHandler(cs CacheService) CacheHandler {
	return &cacheHandler{
		CacheService: cs,
	}
}

// PutItem handles putting an item into cache and returns generated uuid to client.
func (h *cacheHandler) PutItem(w http.ResponseWriter, r *http.Request) {
	item, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}

	key, err := h.CacheService.SetAndGenerateUUID(item)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}

	payload := map[string]string{"key": key}
	respondJSON(w, http.StatusOK, payload)
}

// PutItemWithKey handles putting an item into cache with key specified by client.
func (h *cacheHandler) PutItemWithKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	item, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}

	if err := h.CacheService.Set(key, item); err != nil {
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}

	respondNoContent(w)
}

// GetItem handles getting an item from cache.
func (h *cacheHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	item, err := h.CacheService.Get(key)
	if err != nil {
		if xerrors.Is(err, bigcache.ErrEntryNotFound) {
			respondJSON(w, http.StatusNotFound, buildErrorPayload(err))
			return
		}
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}

	respondJSON(w, http.StatusOK, item)
}

// DeleteItem handles deleting an item from cache by key.
func (h *cacheHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if err := h.CacheService.Delete(key); err != nil {
		if xerrors.Is(err, bigcache.ErrEntryNotFound) {
			respondJSON(w, http.StatusNotFound, buildErrorPayload(err))
			return
		}
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}

	respondNoContent(w)
}

// ListItem handles getting all items in cache.
func (h *cacheHandler) ListItem(w http.ResponseWriter, r *http.Request) {
	items, err := h.CacheService.List()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}

	respondJSON(w, http.StatusOK, items)
}

// GetStats handles getting cache stats.
func (h *cacheHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.CacheService.Stats())
}

// GetStatus handles getting cache status (length and capacity).
func (h *cacheHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	payload := map[string]interface{}{
		"length":   h.CacheService.Len(),
		"capacity": h.CacheService.Capacity(),
	}
	respondJSON(w, http.StatusOK, payload)
}

func buildErrorPayload(err error) map[string]string {
	return map[string]string{"message": err.Error()}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&payload); err != nil {
		respondJSON(w, http.StatusInternalServerError, buildErrorPayload(err))
		return
	}
}

func respondNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
