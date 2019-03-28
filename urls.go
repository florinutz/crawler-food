package kvstore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getSource(req http.Request, client *http.Client) ([]byte, error) {
	resp, err := client.Do(&req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	html = bytes.TrimSuffix(html, []byte{10})

	return html, nil
}

type kvWithErr struct {
	KV
	err error
}

type UrlFetchedCallback func(string, []byte, error)

// FetchUrls loads new data from http
func FetchUrls(requests []http.Request, generalTimeout time.Duration, client *http.Client, gotUrl UrlFetchedCallback) (
	store []KV, errs []error) {

	c := make(chan kvWithErr, len(requests))

	if client == nil {
		client = getDefaultClient(client)
	}

	for _, u := range requests {
		go fetchAsync(u, c, client)
	}

	if generalTimeout == 0 {
		generalTimeout = time.Duration(len(requests)*5) * time.Second
	}

	for i := 0; i < len(requests); i++ {
		select {
		case kve := <-c:
			if kve.err != nil {
				errs = append(errs, kve.err)
				continue
			}
			store = set(kve.Key, kve.Value, store)
			if gotUrl != nil {
				gotUrl(string(kve.Key), kve.Value, kve.err)
			}
		case <-time.After(generalTimeout):
			errs = append(errs, fmt.Errorf("generalTimeout after %s", generalTimeout))
		}
	}

	return
}

func getDefaultClient(client *http.Client) *http.Client {
	client = http.DefaultClient
	client.Timeout = 20 * time.Second
	return client
}

func fetchAsync(req http.Request, output chan<- kvWithErr, client *http.Client) {
	html, err := getSource(req, client)
	output <- kvWithErr{
		KV: KV{
			Key:   []byte(req.URL.String()),
			Value: html,
		},
		err: err,
	}
}

// Set or add a value
func set(key []byte, value []byte, store []KV) (newStore []KV) {
	newStore = make([]KV, len(store))

	for i, existing := range store {
		if string(key) == string(existing.Key) {
			newStore[i].Value = value
			return
		}
	}

	return append(store, KV{Key: key, Value: value})
}

func get(key []byte, store []KV) *KV {
	for _, kv := range store {
		if string(key) == string(kv.Key) {
			return &kv
		}
	}

	return nil
}
