package url_mock

var Urls []*Url

// Url - not using a map here so I can save this as json
type Url struct {
	Url     string `json:"url"`
	Content []byte `json:"content"`
}
