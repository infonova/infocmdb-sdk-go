package cmdb

func (i *InfoCMDB) Query(query string, out interface {}, params map[string]string) (err error) {
	if err = i.Login(); err != nil {
		return
	}

	return i.Client.Put("/apiV2/query", &out, params)
}
