package cacheapi

import (
	"bytes"
	"encoding/json"

	"github.com/allegro/bigcache"
	"github.com/google/uuid"
)

// Item represents cache entry with key-value.
type Item map[string]interface{}

// CacheService provides function for cache.
type CacheService interface {
	Set(key string, v []byte) error
	SetAndGenerateUUID(v []byte) (string, error)
	Get(key string) (Item, error)
	Delete(key string) error
	List() ([]Item, error)
	Stats() bigcache.Stats
	Len() int
	Capacity() int
}

type cacheService struct {
	*bigcache.BigCache
}

// NewCacheService returns a new CacheService.
func NewCacheService(config bigcache.Config) CacheService {
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		panic(err)
	}
	return &cacheService{
		BigCache: cache,
	}
}

// SetAndGenerateUUID wraps BigCache.Set and returns generated uuid.
func (s *cacheService) SetAndGenerateUUID(v []byte) (string, error) {
	id := uuid.New().String()
	if err := s.BigCache.Set(id, v); err != nil {
		return "", err
	}
	return id, nil
}

// Get wraps BigCache.Get with decoding.
func (s *cacheService) Get(key string) (Item, error) {
	b, err := s.BigCache.Get(key)
	if err != nil {
		return nil, err
	}
	var t interface{}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&t); err != nil {
		return nil, err
	}
	return map[string]interface{}{key: t}, nil
}

// List returns All items in cache.
func (s *cacheService) List() ([]Item, error) {
	items := make([]Item, 0, s.BigCache.Len())
	iter := s.BigCache.Iterator()
	for iter.SetNext() {
		v, err := iter.Value()
		if err != nil {
			return nil, err
		}

		var t interface{}
		if err := json.NewDecoder(bytes.NewReader(v.Value())).Decode(&t); err != nil {
			return nil, err
		}
		items = append(items, map[string]interface{}{v.Key(): t})
	}
	return items, nil
}
