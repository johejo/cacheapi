package cacheapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi"

	"github.com/allegro/bigcache"
)

func TestCacheHandler_PutItem(t *testing.T) {
	cs := NewCacheService(bigcache.DefaultConfig(-1))
	ch := NewCacheHandler(cs)

	mux := http.NewServeMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/", strings.NewReader(`{"foo": "bar"}`))

	mux.HandleFunc("/", ch.PutItem)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("invalid status code: code=%d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Errorf("failed to decode body: %s", err)
	}
	if _, ok := body["key"].(string); !ok {
		t.Errorf("failed to get key from body: %v", body)
	}
}

func TestCacheHandler_PutItemWithKey(t *testing.T) {
	cs := NewCacheService(bigcache.DefaultConfig(-1))
	ch := NewCacheHandler(cs)

	mux := chi.NewMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/kkk", strings.NewReader(`{"foo": "var"}`))

	mux.HandleFunc("/{key}", ch.PutItemWithKey)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("invalid status code: code=%d", rec.Code)
	}
}

func TestCacheHandler_GetItem(t *testing.T) {
	cs := NewCacheService(bigcache.DefaultConfig(-1))
	ch := NewCacheHandler(cs)

	mux := chi.NewMux()

	// 200 OK
	{
		b := []byte(`{"foo": "bar"}`)
		const key = "kkk"
		if err := cs.Set(key, b); err != nil {
			t.Errorf("failed to set value to cache: %s", err)
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/"+key, nil)

		mux.HandleFunc("/{key}", ch.GetItem)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("invalid status code: %d", rec.Code)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode body: %s", err)
		}

		expected := map[string]interface{}{
			key: map[string]interface{}{
				"foo": "bar"},
		}
		if !reflect.DeepEqual(body, expected) {
			t.Errorf("failed to compare body: actual=%v. expected=%v", body, expected)
		}
	}

	// 404 Not found
	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/zzz", nil)

		mux.HandleFunc("/{key}", ch.GetItem)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("invalid status code: %d", rec.Code)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode body: %s", err)
		}

		expected := map[string]interface{}{
			"message": "Entry not found",
		}
		if !reflect.DeepEqual(body, expected) {
			t.Errorf("failed to compare body: actual=%v. expected=%v", body, expected)
		}
	}
}

func TestCacheHandler_ListItem(t *testing.T) {
	cs := NewCacheService(bigcache.DefaultConfig(-1))
	ch := NewCacheHandler(cs)

	bs := [][]byte{
		[]byte(`{"foo": "bar"}`),
		[]byte(`{"foo2": "bar2"}`),
	}
	for i, b := range bs {
		if err := cs.Set(strconv.Itoa(i), b); err != nil {
			t.Errorf("failed to set: %s", err)
		}
	}

	mux := http.NewServeMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux.HandleFunc("/", ch.ListItem)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("invalid status code: %d", rec.Code)
	}
	var body []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Errorf("failed to decode body: %s", err)
	}

	expected := []map[string]interface{}{
		{
			"0": map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			"1": map[string]interface{}{
				"foo2": "bar2",
			},
		},
	}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("failed to compare body: actual=%v. expected=%v", body, expected)
	}
}

func TestCacheHandler_DeleteItem(t *testing.T) {
	cs := NewCacheService(bigcache.DefaultConfig(-1))
	ch := NewCacheHandler(cs)

	mux := chi.NewMux()

	// 200 OK
	{
		b := []byte(`{"foo": "bar"}`)
		const key = "kkk"
		if err := cs.Set(key, b); err != nil {
			t.Errorf("failed to set value to cache: %s", err)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/"+key, nil)

		mux.HandleFunc("/{key}", ch.DeleteItem)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("invalid status code: %d", rec.Code)
		}
	}

	// 404 Not found
	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/zzz", nil)

		mux.HandleFunc("/{key}", ch.GetItem)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("invalid status code: %d", rec.Code)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode body: %s", err)
		}

		expected := map[string]interface{}{
			"message": "Entry not found",
		}
		if !reflect.DeepEqual(body, expected) {
			t.Errorf("failed to compare body: actual=%v. expected=%v", body, expected)
		}
	}
}

func TestCacheHandler_GetStats(t *testing.T) {
	cs := NewCacheService(bigcache.DefaultConfig(-1))
	ch := NewCacheHandler(cs)

	mux := http.NewServeMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux.HandleFunc("/", ch.GetStats)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("invalid status code: %d", rec.Code)
	}

	var actual map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&actual); err != nil {
		t.Errorf("failed to decode body: %s", err)
	}

	hists := actual["hits"]
	misses := actual["misses"]
	deleteHits := actual["delete_hits"]
	deleteMisses := actual["delete_misses"]
	collisions := actual["collisions"]
	if hists == nil || misses == nil || deleteHits == nil || deleteMisses == nil || collisions == nil {
		t.Errorf("misssing key: %v", actual)
	}
}

func TestCacheHandler_GetStatus(t *testing.T) {
	cs := NewCacheService(bigcache.DefaultConfig(-1))
	ch := NewCacheHandler(cs)

	mux := http.NewServeMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux.HandleFunc("/", ch.GetStatus)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("invalid status code: %d", rec.Code)
	}

	var actual map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&actual); err != nil {
		t.Errorf("failed to decode body: %s", err)
	}

	if _, ok := actual["length"]; !ok {
		t.Errorf(`failed to get key "length": %v`, actual)
	}
	if _, ok := actual["capacity"]; !ok {
		t.Errorf(`failed to get key "capacity": %v`, actual)
	}
}
