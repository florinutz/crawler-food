package kvstore

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	var responseContent string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, responseContent)
	}))

	defer ts.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore(tt.args.client, tt.args.timeout, tt.args.urlFetched)
			kvs := []KV{
				{
					Key:   []byte("keyOne"),
					Value: []byte("aValue"),
				},
				{
					Key:   []byte("keyTwo"),
					Value: []byte("anotherValue"),
				},
			}

			buf := &bytes.Buffer{}

			store.Write(kvs, buf)

			responseContent = buf.String()

			fetchedKvs, errs := store.Fetch([]string{ts.URL})
			if len(errs) == 0 {
				t.Fatal("fetch errors")
			}

			if !reflect.DeepEqual(fetchedKvs, kvs) {
				t.Fatal("mismatch")
			}
		})
	}
}
