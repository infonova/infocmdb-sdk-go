package cmdb

func (i *InfoCMDB) Query(query string, out interface {}, params map[string]string) (err error) {
	return i.Client.Post("/apiV2/query", &out, params)
}