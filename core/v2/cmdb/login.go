package cmdb

import (
	"errors"
	"github.com/infonova/infocmdb-lib-go/core/v1/cmdb"
	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb/client"
	log "github.com/sirupsen/logrus"
)

type LoginTokenReturn struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (i *InfoCMDB) Login() (err error) {

	if i.Config.ApiKey != "" {
		log.Debug("already logged in")
		return nil
	}
	if i.Config.Username == "" || i.Config.Password == "" {
		return cmdb.ErrNoCredentials
	}

	var loginResult LoginTokenReturn
	params := map[string]string{
		"username": i.Config.Username,
		"password": i.Config.Password,
		"lifetime": "600",
	}

	if i.Client == nil {
		i.Client = client.NewClient(i.Config.Url)
	}

	i.Client.SetHostURL(i.Config.Url)

	err = i.Client.PostForm("/apiV2/auth/token", &loginResult, params)
	if err != nil {
		err = errors.New("Failed to fetch token: " + err.Error())
		return
	}

	if loginResult.Data.Token == "" {
		return cmdb.ErrLoginFailed
	}

	i.Config.ApiKey = loginResult.Data.Token
	i.Client.SetAuthToken(i.Config.ApiKey)
	return
}
