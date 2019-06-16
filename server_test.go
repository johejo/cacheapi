package cacheapi

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
)

func TestServer_Run(t *testing.T) {
	api := NewCacheAPI(newDummyCacheHandler())
	const (
		host = ""
		port = 9999
	)
	s := NewServer(api, host, port)
	ch := make(chan bool, 1)
	go func() {
		ch <- true
		log.Fatal(s.Run())
	}()

	<-ch
	addr := fmt.Sprintf("http://127.0.0.1:%d", port)
	_, err := http.Get(addr)
	if err != nil {
		t.Errorf("failed to send request: addr=%s, err=%s", addr, err)
	}

	// should be net.Error
	addr = fmt.Sprintf("http://127.0.0.1:%d", 1111)
	_, err = http.Get(addr)
	if err != nil {
		if _, ok := err.(net.Error); !ok {
			log.Fatal(err)
		}
	}
}
