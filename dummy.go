package cacheapi

import (
	"net/http"

	"github.com/allegro/bigcache"
)

type dummyCacheService struct{}

func newDummyCacheService() CacheService { return &dummyCacheService{} }

func (s *dummyCacheService) SetAndGenerateUUID(v []byte) (string, error) { return "", nil }
func (s *dummyCacheService) Get(key string) (Item, error)                { return nil, nil }
func (s *dummyCacheService) List() ([]Item, error)                       { return nil, nil }
func (s *dummyCacheService) Delete(key string) error                     { return nil }
func (s *dummyCacheService) Set(key string, v []byte) error              { return nil }
func (s *dummyCacheService) Stats() bigcache.Stats                       { return bigcache.Stats{} }
func (s *dummyCacheService) Len() int                                    { return 0 }
func (s *dummyCacheService) Capacity() int                               { return 0 }

type dummyCacheHandler struct{}

func newDummyCacheHandler() CacheHandler { return &dummyCacheHandler{} }

func (h *dummyCacheHandler) PutItem(w http.ResponseWriter, r *http.Request)        {}
func (h *dummyCacheHandler) PutItemWithKey(w http.ResponseWriter, r *http.Request) {}
func (h *dummyCacheHandler) GetItem(w http.ResponseWriter, r *http.Request)        {}
func (h *dummyCacheHandler) DeleteItem(w http.ResponseWriter, r *http.Request)     {}
func (h *dummyCacheHandler) ListItem(w http.ResponseWriter, r *http.Request)       {}
func (h *dummyCacheHandler) GetStats(w http.ResponseWriter, r *http.Request)       {}
func (h *dummyCacheHandler) GetStatus(w http.ResponseWriter, r *http.Request)      {}
