package cacheapi

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCacheAPI(t *testing.T) {
	h := NewCacheHandler(newDummyCacheService())
	api := NewCacheAPI(h)
	ts := httptest.NewServer(api.Handler())

	var (
		resp     *http.Response
		respBody string
	)

	// ListItem
	{
		resp, respBody = testRequest(t, ts, "GET", cachePath, "", nil)
		if resp.StatusCode != http.StatusOK && respBody == "[]" {
			t.Errorf("failed to get item list: %s", respBody)
		}
	}

	// GetItem
	{
		resp, respBody = testRequest(t, ts, "GET", cachePath+"/zzz", "", nil)
		if resp.StatusCode != http.StatusNotFound && resp.Header.Get("Content-Type") != "application/json" {
			t.Errorf("should be 404: %d", resp.StatusCode)
		}
	}

	// DeleteItem
	{
		resp, respBody = testRequest(t, ts, "DELETE", cachePath+"/zzz", "", nil)
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("should be 204: %d", resp.StatusCode)
		}
	}

	// PutItem
	{
		resp, respBody = testRequest(t, ts, "PUT", cachePath, "application/json", strings.NewReader(`{"foo": "bar"}`))
		if resp.StatusCode != http.StatusOK {
			t.Errorf("should be 200: %d", resp.StatusCode)
		}
	}

	// PutItemWithKey
	{
		resp, respBody = testRequest(t, ts, "PUT", cachePath+"/zzz", "application/json", strings.NewReader(`{"foo": "bar"}`))
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("should be 204: %d", resp.StatusCode)
		}
	}

	// GetStatus
	{
		resp, respBody = testRequest(t, ts, "GET", cachePath+"/status", "", nil)
		if resp.StatusCode != http.StatusNotFound && resp.Header.Get("Content-Type") != "application/json" {
			t.Errorf("should be 404: %d", resp.StatusCode)
		}
	}

	// GetStats
	{
		resp, respBody = testRequest(t, ts, "GET", basePath+"/stats", "", nil)
		if resp.StatusCode != http.StatusNotFound && resp.Header.Get("Content-Type") != "application/json" {
			t.Errorf("should be 404: %d", resp.StatusCode)
		}
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, contentType string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
