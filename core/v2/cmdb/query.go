package cmdb

import "gopkg.in/resty.v1"

type queryParams struct {
	Params map[string]string `json:"params"`
}

type queryRequest struct {
	Query queryParams `json:"query"`
}

func (i *InfoCMDB) Query(query string, out interface{}, params map[string]string) (resp *resty.Response, err error) {
	if err = i.Login(); err != nil {
		return
	}

	r := queryRequest{
		Query: queryParams{
			Params: params,
		},
	}

	resp, err = i.Client.NewRequest().
		SetResult(out).
		SetBody(r).
		SetAuthToken(i.Config.ApiKey).
		Put("/apiV2/query/execute/" + query)

	return
}
