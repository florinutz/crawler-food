package kvstore

import (
	"io"
)

// KV - not using a map here so I can save this as json
type KV struct {
	Key   []byte `json:"k"`
	Value []byte `json:"v"`
}

type ValueScrambler interface {
	EncodeValue(value []byte) (encoded []byte)
	DecodeValue(encoded []byte) (decoded []byte, err error)
}

type KeyScrambler interface {
	EncodeKey(value []byte) (encoded []byte)
	DecodeKey(encoded []byte) (decoded []byte, err error)
}

type StoreScrambler interface {
	EncodeStore([]KV) (encoded []byte, err error)
}

type Scrambler interface {
	KeyScrambler
	ValueScrambler
	StoreScrambler
}

type IO interface {
	Read(from io.Reader) (kvs []KV, err error)
	Write(kv []KV, to io.Writer) (err error)
}

type DataLoader interface {
	Fetch(keys []string) (set []KV, errs []error)
}

type KVStore interface {
	DataLoader
	Scrambler
	IO
}
