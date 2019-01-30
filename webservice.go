package infocmdblibrary

import (
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Webservice struct {
	client *CmdbWebClient
}

// Webservice queries a given webservice with all params supplied
// Returns err != nil if query fails
func (w *Webservice) Webservice(ws string, params url.Values) (string, error) {
	log.Debugf("Webservice: %s, Params: %v", ws, params)
	r, err := w.client.Post("query", ws, params)
	if err != nil {
		return "", err
	}
	return r, nil
}
