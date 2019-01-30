package infocmdblibrary

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
)

type Webservice struct {
	client *CmdbWebClient
}

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("INFOCMDB_WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

// Webservice queries a given webservice with all params supplied
// Returns err != nil if query fails
func (w *Webservice) Webservice(ws string, params url.Values) (string, error) {
	log.Debugf("Webservice: %s, Params: %v\n", ws, params)
	r, err := w.client.Post("query", ws, params)
	if err != nil {
		return "", err
	}
	return r, nil
}
