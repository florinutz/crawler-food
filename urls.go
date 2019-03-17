package kv_store

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func getSource(visitUrl string) ([]byte, error) {
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

// FetchUrls loads new data
func FetchUrls(wantedUrls []string, timeout time.Duration) (errs []error) {
	count := len(wantedUrls)

	c := make(chan kvWithErr, count)

	for _, u := range wantedUrls {
		go fetchAsync(u, c)
	}

	var store []KV

	for i := 0; i < count; i++ {
		select {
		case block := <-c:
			if block.err != nil {
				errs = append(errs, block.err)
				continue
			}
			set(block.Key, block.Value, store)
			fmt.Printf("* loaded %s\n", block.Key)
		case <-time.After(timeout):
			errs = append(errs, fmt.Errorf("timeout after %s", timeout))
		}
	}

	return
}

func fetchAsync(url string, blockChan chan<- kvWithErr) {
	html, err := getSource(url)
	blockChan <- kvWithErr{
		KV: KV{
			Key:   url,
			Value: html,
		},
		err: err,
	}
}

// Set or add a value
func set(key string, value []byte, store []KV) {
	for i, existing := range store {
		if key == existing.Key {
			store[i].Value = value
			return
		}
	}
	store = append(store, KV{Key: key, Value: value})
}
