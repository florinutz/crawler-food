package kvstore

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// todo modify the tests
// todo update readme

func getSource(visitUrl string, transport *http.Client) ([]byte, error) {
	_, err := url.Parse(visitUrl)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(visitUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return html, nil
}

type kvWithErr struct {
	KV
	err error
}

type UrlFetchedCallback func(string, []byte, error)

// FetchUrls loads new data from http
func FetchUrls(wantedUrls []string, generalTimeout time.Duration, client *http.Client, gotUrl UrlFetchedCallback) (
	store []KV, errs []error) {
	count := len(wantedUrls)

	c := make(chan kvWithErr, count)

	if client == nil {
		client = getDefaultClient(client)
	}

	for _, u := range wantedUrls {
		go fetchAsync(u, c, client)
	}

	for i := 0; i < count; i++ {
		select {
		case block := <-c:
			if block.err != nil {
				errs = append(errs, block.err)
				continue
			}
			store = set(block.Key, block.Value, store)
			if gotUrl != nil {
				gotUrl(block.Key, block.Value, block.err)
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

func fetchAsync(url string, blockChan chan<- kvWithErr, client *http.Client) {
	html, err := getSource(url, client)
	blockChan <- kvWithErr{
		KV: KV{
			Key:   url,
			Value: html,
		},
		err: err,
	}
}

// Set or add a value
func set(key string, value []byte, store []KV) (newStore []KV) {
	for i, existing := range store {
		if key == existing.Key {
			newStore[i].Value = value
			return
		}
	}

	return append(store, KV{Key: key, Value: value})
}
