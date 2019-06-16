package cacheapi

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	basePath   = "/api/v0"
	cachePath  = basePath + "/cache"
	keyPath    = cachePath + "/{key}"
	statsPath  = basePath + "/stats"
	statusPath = cachePath + "/status"
)

// CacheAPI abstracts cache for http.Handler.
type CacheAPI interface {
	Handler() http.Handler
}

type cacheAPI struct {
	router chi.Router
	CacheHandler
}

// NewCacheAPI returns a new CacheAPI.
func NewCacheAPI(ch CacheHandler) CacheAPI {
	return &cacheAPI{
		router:       chi.NewRouter(),
		CacheHandler: ch,
	}
}

// Handler returns routed http.Handler.
func (api *cacheAPI) Handler() http.Handler {
	if GetEnv("REQUEST_LOGGING", "off") == "on" {
		api.router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: GetLogger()}))
	}

	api.router.Use(middleware.AllowContentType("application/json"))

	api.router.Get(cachePath, api.CacheHandler.ListItem)
	api.router.Put(cachePath, api.CacheHandler.PutItem)
	api.router.Put(keyPath, api.CacheHandler.PutItemWithKey)
	api.router.Get(keyPath, api.CacheHandler.GetItem)
	api.router.Delete(keyPath, api.CacheHandler.DeleteItem)
	api.router.Get(statusPath, api.CacheHandler.GetStatus)
	api.router.Get(statsPath, api.CacheHandler.GetStats)

	return api.router
}
