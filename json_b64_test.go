package kvstore

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadWrite(t *testing.T) {
	type readerArgs struct {
		sliceDecoder   Decoder
		contentDecoder ContentDecoder
	}
	type writerArgs struct {
		kvs            []KV
		sliceEncoder   Encoder
		contentEncoder ContentEncoder
	}

	kvs := []KV{
		{"https://something.com/some-resource.html", []byte("some content")},
		{"https://something-else.com/some-other-resource.html", []byte("some other content")},
	}

	tests := []struct {
		name          string
		writerArgs    writerArgs
		readerArgs    readerArgs
		wantMismatch  bool
		wantWriterErr bool
		wantReaderErr bool
	}{
		{
			name: "full successful cycle",
			writerArgs: writerArgs{
				kvs:            kvs,
				sliceEncoder:   DefaultEncoder,
				contentEncoder: DefaultContentEncoder,
			},
			readerArgs: readerArgs{
				sliceDecoder:   DefaultDecoder,
				contentDecoder: DefaultContentDecoder,
			},
		},
		{
			name: "encoding mismatch",
			writerArgs: writerArgs{
				kvs:            kvs,
				sliceEncoder:   DefaultEncoder,
				contentEncoder: DefaultContentEncoder,
			},
			readerArgs: readerArgs{
				sliceDecoder:   DefaultDecoder,
				contentDecoder: nil,
			},
			wantMismatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer

			originalKvs := make([]KV, len(kvs))
			copy(originalKvs, kvs)

			if err := Write(kvs, &buffer, tt.writerArgs.sliceEncoder, tt.writerArgs.contentEncoder); (err != nil) != tt.wantWriterErr {
				t.Errorf("Write() error = %v, wantWriterErr %v", err, tt.wantWriterErr)
				return
			}

			if !reflect.DeepEqual(kvs, originalKvs) {
				t.Error("write had side effects (it changed the kvs)")
				return
			}

			got, err := Read(&buffer, tt.readerArgs.sliceDecoder, tt.readerArgs.contentDecoder)
			if (err != nil) != tt.wantReaderErr {
				t.Errorf("Read() error = %v, wantReaderErr %v", err, tt.wantReaderErr)
				return
			}

			if !tt.wantMismatch && !reflect.DeepEqual(kvs, got) {
				t.Errorf("data was corrupted on the way: encoded \n%q\n, got \n%q", kvs, got)
			}
		})
	}
}
