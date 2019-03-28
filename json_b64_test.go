package kvstore

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	type args struct {
		client     *http.Client
		timeout    time.Duration
		urlFetched UrlFetchedCallback
	}
	tests := []struct {
		name string
		args args
		want *JsonB64UrlStore
	}{
		{
			name: "simple",
		},
	}

	var server1Reply, server2Reply string

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, server1Reply)
	}))
	defer server1.Close()
	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, server2Reply)
	}))
	defer server2.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server1Reply = "ana are mere"
			server2Reply = "ciresel vine si cere"

			store := NewStore(tt.args.client, tt.args.timeout, tt.args.urlFetched)
			kvs := []KV{
				{
					Key:   []byte(server1.URL),
					Value: []byte(server1Reply),
				},
				{
					Key:   []byte(server2.URL),
					Value: []byte(server2Reply),
				},
			}

			buf := &bytes.Buffer{}

			err := store.Write(kvs, buf)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Printf("%s\n", buf)

			readKvs, err := store.Read(buf)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(readKvs, kvs) {
				t.Fatal("mismatch after write/read")
			}

			req1, _ := http.NewRequest("GET", server1.URL, nil)
			req2, _ := http.NewRequest("GET", server2.URL, nil)
			reqs := []http.Request{*req1, *req2}

			fetchedKvs, errs := store.Fetch(reqs)
			if len(errs) > 0 {
				var aux []string
				for _, err := range errs {
					aux = append(aux, err.Error())
				}
				t.Fatalf("fetch errors: \n* %s", strings.Join(aux, "\n * "))
			}

			// can't do reflect.DeepEquals since order is lost
			for _, kv := range kvs {
				if get(kv.Key, fetchedKvs) == nil {
					t.Fatalf("url '%s' not fetched", string(kv.Key))
				}
			}
		})
	}
}
