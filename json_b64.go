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

type JsonB64UrlStore struct {
	timeout time.Duration
	client  *http.Client
	gotUrl  UrlFetchedCallback
}

func NewStore(client *http.Client, timeout time.Duration, urlFetched UrlFetchedCallback) *JsonB64UrlStore {
	return &JsonB64UrlStore{timeout: timeout, client: client, gotUrl: urlFetched}
}

func (store *JsonB64UrlStore) EncodeStore(kvs []KV) (encoded []byte, err error) {
	// don't tamper with the source
	var aux []KV
	for _, u := range kvs {
		newKey := store.EncodeKey([]byte(u.Key))
		newValue := store.EncodeValue(u.Value)
		aux = append(aux, KV{
			Key:   newKey,
			Value: newValue,
		})
	}

	return json.Marshal(aux)
}

func (store *JsonB64UrlStore) Read(from io.Reader) (kvs []KV, err error) {
	var tmp []KV

	bytes, err := ioutil.ReadAll(from)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &tmp)

	for _, kv := range tmp {
		key, err := store.DecodeKey(kv.Key)
		if err != nil {
			return nil, err
		}
		val, err := store.DecodeValue(kv.Value)
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, KV{
			Key:   key,
			Value: val,
		})
	}
	return
}

func (store *JsonB64UrlStore) Write(kv []KV, to io.Writer) (err error) {
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

func (store *JsonB64UrlStore) Fetch(keys []string) (set []KV, errs []error) {
	return FetchUrls(keys, store.timeout, store.client, store.gotUrl)
}

func (store *JsonB64UrlStore) EncodeKey(value []byte) (encoded []byte) {
	return encodeToB64(value)
}

func (store *JsonB64UrlStore) DecodeKey(encoded []byte) (decoded []byte, err error) {
	return decodeB64(encoded)
}

func (store *JsonB64UrlStore) EncodeValue(value []byte) (encoded []byte) {
	return encodeToB64(value)
}

func (store *JsonB64UrlStore) DecodeValue(encoded []byte) (decoded []byte, err error) {
	return decodeB64(encoded)
}

func decodeB64(encoded []byte) (decoded []byte, err error) {
	decoded = make([]byte, base64.RawStdEncoding.DecodedLen(len(encoded)))
	_, err = base64.RawStdEncoding.Decode(decoded, encoded)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode base64")
	}
	return
}

func encodeToB64(value []byte) (encoded []byte) {
	encoded = make([]byte, base64.RawStdEncoding.EncodedLen(len(value)))
	base64.RawStdEncoding.Encode(encoded, value)
	return
}
