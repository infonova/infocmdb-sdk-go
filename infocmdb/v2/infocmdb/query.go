package infocmdb

import (
	"github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb/client"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type queryParams struct {
	Params map[string]string `json:"params"`
}

type queryRequest struct {
	Query queryParams `json:"query"`
}

func (cmdb *Cmdb) Query(query string, out interface{}, params map[string]string) (err error) {
	log.Debugf("Querying webservice %v with params %v", query, params)

	if err = cmdb.Login(); err != nil {
		return
	}

	body := queryRequest{
		Query: queryParams{
			Params: params,
		},
	}

	var respError client.ResponseError

	resp, err := cmdb.Client.Execute(resty.MethodPut, "/apiV2/query/execute/"+query,
		func(request *resty.Request) *resty.Request {
			return request.
				SetResult(out).
				SetBody(body).
				SetError(&respError)
		})

	if err != nil {
		return
	}

	if resp.IsError() {
		log.Debugf("Status: %v, Error result: %v", resp.StatusCode(), respError)
		return respError
	}

	log.Debugf("Response: %s", resp.String())
	log.Debugf("Mapped result: %+v", out)
	return
}

func (cmdb *Cmdb) QueryRaw(query string, params map[string]string) (r string, err error) {
	if err = cmdb.Login(); err != nil {
		return
	}

	qr := queryRequest{
		Query: queryParams{
			Params: params,
		},
	}

	var respError client.ResponseError
	resp, err := cmdb.Client.Execute(resty.MethodPut, "/apiV2/query/execute/"+query,
		func(request *resty.Request) *resty.Request {
			return request.
				SetBody(qr).
				SetError(&respError)
		})

	if resp == nil {
		return
	}

	if resp.IsError() {
		return "", respError
	}

	return resp.String(), nil
}
