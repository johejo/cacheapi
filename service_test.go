package cacheapi

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/allegro/bigcache"
)

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func TestCacheService_SetAndGenerateUUID(t *testing.T) {
	s := NewCacheService(bigcache.DefaultConfig(-1))
	b := []byte("foo")

	key, err := s.SetAndGenerateUUID(b)
	if err != nil {
		t.Errorf("failed to set and generate uuid: %s", err)
	}
	if !isValidUUID(key) {
		t.Errorf("failed to generate valid uuid: %s", err)
	}
}

func TestCacheService_Get(t *testing.T) {
	var (
		actual, expected Item
		b                []byte
		err              error
	)
	s := NewCacheService(bigcache.DefaultConfig(-1))
	b = []byte(`{"key": "value"}`)

	key, err := s.SetAndGenerateUUID(b)
	if err != nil {
		t.Errorf("failed to set and generate uuid: %s", err)
	}

	actual, err = s.Get(key)
	if err != nil {
		t.Errorf("failed to get: %s", err)
	}

	expected = map[string]interface{}{"key": "value"}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("faield to decode item: actual=%v, expected=%v", actual, expected)
	}
}

func TestCacheService_List(t *testing.T) {
	var (
		actual, expected []Item
		bs               [][]byte
		err              error
	)
	s := NewCacheService(bigcache.DefaultConfig(-1))
	bs = [][]byte{
		[]byte(`{"key": "value"}`),
		[]byte(`{"key2": "value2"}`),
	}
	for _, b := range bs {
		if err := s.Set(uuid.New().String(), b); err != nil {
			t.Errorf("failed to set: %s", err)
		}
	}

	actual, err = s.List()
	if err != nil {
		t.Errorf("failed to get list: %s", err)
	}

	expected = []Item{
		map[string]interface{}{"key": "value"},
		map[string]interface{}{"key2": "value2"},
	}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("failed to decode item list: actual=%s, expected=%s", actual, expected)
	}
}
