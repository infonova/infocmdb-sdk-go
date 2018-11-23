package infocmdbGoLib

import (
	"net/url"
)

type Webservice struct {
	client *cmdbWebClient
}

// Webservice queries a given webservice with all params supplied
// Returns err != nil if query fails
func (w *Webservice) Webservice(ws string, params url.Values) (string, error) {

	r, err := w.client.Post("query", ws, params)
	if err != nil {
		return "", err
	}
	return r, nil
}
