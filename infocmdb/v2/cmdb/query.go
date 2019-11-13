package cmdb

import (
	"github.com/infonova/infocmdb-sdk-go/infocmdb/v2/cmdb/client"
)

type queryParams struct {
	Params map[string]string `json:"params"`
}

type queryRequest struct {
	Query queryParams `json:"query"`
}

func (i *Cmdb) Query(query string, out interface{}, params map[string]string) (err error) {
	if err = i.Login(); err != nil {
		return
	}

	r := queryRequest{
		Query: queryParams{
			Params: params,
		},
	}

	var respError client.ResponseError

	resp, err := i.Client.NewRequest().
		SetResult(out).
		SetBody(r).
		SetAuthToken(i.Config.ApiKey).
		SetError(&respError).
		Put("/apiV2/query/execute/" + query)

	if resp != nil && resp.IsError() {
		return respError
	}

	return
}

func (i *Cmdb) QueryRaw(query string, params map[string]string) (r string, err error) {
	if err = i.Login(); err != nil {
		return
	}

	qr := queryRequest{
		Query: queryParams{
			Params: params,
		},
	}

	var respError client.ResponseError
	resp, err := i.Client.NewRequest().
		SetBody(qr).
		SetAuthToken(i.Config.ApiKey).
		SetError(&respError).
		Put("/apiV2/query/execute/" + query)

	if resp == nil {
		return
	}

	if resp.IsError() {
		return "", respError
	}

	return resp.String(), nil
}
