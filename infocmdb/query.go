package infocmdb

import (
	log "github.com/sirupsen/logrus"
)

// QueryWebservices allows you to call a generic webservice(arg1: ws) with the providing params
// Return: json string
func (c *Client) QueryWebservice(ws string, params map[string]string) (resp string, err error) {
	log.Debugf("Querying webservice %v with params %v", ws, params)

	if err = c.v2.Login(); err != nil {
		return
	}

	resp, err = c.v2.QueryRaw(ws, params)
	if err != nil {
		log.Error("Error: ", err)
	}

	log.Debugf("Result: %v", resp)
	return
}

// Query allows you to call a generic webservice(arg1: ws) with the providing params and a reference
// to a result. It will take the built in resty function to deserialize the result
// Return: error
func (c *Client) Query(ws string, out interface{}, params map[string]string) (err error) {
	log.Debugf("Querying webservice %v with params %v", ws, params)

	if err = c.v2.Login(); err != nil {
		return
	}

	if err = c.v2.Query(ws, out, params); err != nil {
		log.Error("Error: ", err)
	}

	log.Debugf("Result: %v", out)
	return
}
