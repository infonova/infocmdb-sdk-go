package cmdb

import (
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
		return ErrNoCredentials
	}

	var loginResult LoginTokenReturn
	params := map[string]string{
		"username": i.Config.Username,
		"password": i.Config.Password ,
		"lifetime": "600",
	}

	i.Client.SetHostURL(i.Config.Url)

	err = i.Client.Post("/apiV2/auth/token",&loginResult,params)
	if err != nil {
		return
	}

	if loginResult.Data.Token == "" {
		return ErrLoginFailed
	}

	i.Config.ApiKey = loginResult.Data.Token
	i.Client.SetAuthToken(i.Config.ApiKey)
	return
}