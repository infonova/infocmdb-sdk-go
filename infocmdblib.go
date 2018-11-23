package infocmdblibrary

import (
	"errors"
)

type Credentials struct {
	ApiKey   string
	Username string
	Password string
}

type Config struct {
	Path string
	URL  string
}

type InfoCmdbGoLib struct {
	WS Webservice
	WC cmdbWebClient
}

func NewCMDB(url string, cred Credentials) (InfoCmdbGoLib, error) {
	i := InfoCmdbGoLib{}
	i.WC.apikey = cred.ApiKey
	if cred.ApiKey == "" && cred.Username == "" {
		return i, errors.New("must provide credentials")
	}

	if i.WC.apikey == "" && cred.Username != "" {
		i.Login(url, cred.Username, cred.Password)
	}
	return i, nil
}

func (i *InfoCmdbGoLib) Login(url string, username string, password string) error {
	i.WC = NewCmdbWebClient()
	i.WS.client = &i.WC

	return i.WC.Login(url, username, password)
}
