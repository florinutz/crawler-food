package url_mock

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type Encoder interface {
	Encode([]*Url) []byte
}

type Decoder interface {
	Decode([]byte) ([]*Url, error)
}

type Data struct {
	urls   []*Url
	reader io.Reader
	writer io.Writer
}

// Load reads all urls from data file
func Load(reader io.Reader) (urls []*Url, err error) {
	byteValue, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteValue, &urls)

	return
}

func getReader(filePath string) (io.Reader, error) {
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}

	return jsonFile, nil
}
