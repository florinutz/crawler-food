package kvstore

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type jsonB64UrlStore struct {
	timeout time.Duration
	client  *http.Client
	gotUrl  UrlFetchedCallback
}

func NewStore(timeout time.Duration, client *http.Client, gotUrl UrlFetchedCallback) *jsonB64UrlStore {
	return &jsonB64UrlStore{timeout: timeout, client: client, gotUrl: gotUrl}
}

func (store *jsonB64UrlStore) EncodeStore(kvs []KV) (encoded []byte, err error) {
	// don't tamper with the source
	aux := make([]KV, len(kvs))
	copy(aux, kvs)

	for i, u := range aux {
		newKey := store.EncodeKey([]byte(u.Key))
		newValue := store.EncodeValue(u.Value)
		aux[i].Key = string(newKey)
		aux[i].Value = newValue
	}

	return json.Marshal(aux)
}

func (store *jsonB64UrlStore) Read(from io.Reader) (kvs []KV, err error) {
	bytes, err := ioutil.ReadAll(from)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &kvs)
	return
}

func (store *jsonB64UrlStore) Write(kv []KV, to io.Writer) (err error) {
	b, err := store.EncodeStore(kv)
	if err != nil {
		return err
	}

	n, err := to.Write(b)
	if err != nil {
		return err
	}
	if n < len(b) {
		return io.ErrShortWrite
	}

	return nil
}

func (store *jsonB64UrlStore) Fetch(keys []string) (set []KV, errs []error) {
	return FetchUrls(keys, store.timeout, store.client, store.gotUrl)
}

func (store *jsonB64UrlStore) EncodeKey(value []byte) (encoded []byte) {
	return encodeToB64(value)
}

func (store *jsonB64UrlStore) DecodeKey(encoded []byte) (decoded []byte, err error) {
	return decodeB64(encoded)
}

func (store *jsonB64UrlStore) EncodeValue(value []byte) (encoded []byte) {
	return encodeToB64(value)
}

func (store *jsonB64UrlStore) DecodeValue(encoded []byte) (decoded []byte, err error) {
	return decodeB64(encoded)
}

func decodeB64(encoded []byte) (decoded []byte, err error) {
	decoded = make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	_, err = base64.StdEncoding.Decode(decoded, encoded)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode base64")
	}
	return
}

func encodeToB64(value []byte) (encoded []byte) {
	encoded = make([]byte, base64.StdEncoding.EncodedLen(len(value)))
	base64.StdEncoding.Encode(encoded, value)
	return
}
