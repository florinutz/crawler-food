package kvstore

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// KV - not using a map here so I can save this as json
type KV struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

type ContentEncoder func([]byte) []byte

type ContentDecoder func([]byte) ([]byte, error)

type Encoder func([]KV, ContentEncoder) ([]byte, error)

type Decoder func([]byte, ContentDecoder) ([]KV, error)

var DefaultDecoder Decoder = func(bytes []byte, contentDecoder ContentDecoder) (kvs []KV, err error) {
	err = json.Unmarshal(bytes, &kvs)

	if contentDecoder != nil {
		for i, u := range kvs {
			kvs[i].Value, err = contentDecoder(u.Value)
			if err != nil {
				return nil, errors.Wrapf(err, "error decoding content for '%s'", u.Key)
			}
		}
	}

	return
}

var DefaultEncoder Encoder = func(kvs []KV, contentEncoder ContentEncoder) (bytes []byte, err error) {
	aux := make([]KV, len(kvs))
	copy(aux, kvs)

	if contentEncoder != nil {
		for i, u := range aux {
			aux[i].Value = contentEncoder(u.Value)
		}
	}

	return json.Marshal(aux)
}

var DefaultContentEncoder ContentEncoder = func(bytes []byte) (result []byte) {
	result = make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(result, bytes)
	return
}

var DefaultContentDecoder ContentDecoder = func(encoded []byte) (decoded []byte, err error) {
	decoded = make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	_, err = base64.StdEncoding.Decode(decoded, encoded)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode base64")
	}

	return
}

// Load reads kvs
func Read(reader io.Reader, sliceDecoder Decoder, contentDecoder ContentDecoder) (kvs []KV, err error) {
	if sliceDecoder == nil {
		return nil, errors.New("this has to be decoded somehow")
	}

	byteValue, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	kvs, err = sliceDecoder(byteValue, contentDecoder)

	return
}

// Write writes kvs to writer, while encoding each value and then also sliceEncoding the whole thing
func Write(kvs []KV, writer io.Writer, sliceEncoder Encoder, contentEncoder ContentEncoder) error {
	if sliceEncoder == nil {
		return errors.New("this has to be encoded somehow to bytes")
	}

	b, err := sliceEncoder(kvs, contentEncoder)
	if err != nil {
		return err
	}

	n, err := writer.Write(b)
	if err != nil {
		return err
	}
	if n < len(b) {
		return io.ErrShortWrite
	}

	return nil
}
