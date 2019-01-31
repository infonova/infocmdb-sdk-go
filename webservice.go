package infocmdbclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ErrLoginFailed = errors.New("login status not ok")
)

type ResultLogin struct {
	Status string `json:"status"`
	ApiKey string `json:"apikey"`
}

func (i *InfoCMDB) LoginWithUserPass(url string, username string, password string) error {

	log.Debugf("Opening new Webclient connection. (Url: %s, Username: %s)", url, username)

	reqURL := fmt.Sprintf("%s/api/login/username/%s/password/%s/timeout/21600/method/json", url, username, password)

	resp, err := http.Get(reqURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResult ResultLogin
	err = json.Unmarshal(byteBody, &loginResult)
	if err != nil {
		return err
	}

	if loginResult.Status != "OK" {
		return ErrLoginFailed
	}

	i.Config.ApiKey = loginResult.ApiKey
	return nil
}

func (i *InfoCMDB) LoginWithApiKey(url string, apikey string) error {
	log.Debugf("Opening new Webclient connection using Apikey. (Url: %s, ApiKey: %s)", url, apikey)
	i.Config.ApiUrl = url
	i.Config.ApiKey = apikey
	return nil
}

// Get the api with a given method and parameters as a GetRequest
func (i *InfoCMDB) Get(service string, serviceName string, params url.Values) (string, error) {
	ret, err := i.CallWebservice(http.MethodGet, service, serviceName, params)
	return ret, err
}

// Post the api with a given method and parameters as a GetRequest
func (i *InfoCMDB) Post(service string, serviceName string, params url.Values) (string, error) {
	ret, err := i.CallWebservice(http.MethodPost, service, serviceName, params)
	return ret, err
}

func (i *InfoCMDB) CallWebservice(method string, service string, serviceName string, params url.Values) (response string, err error) {
	if i.Config.ApiKey == "" {
		err = i.Login()
		if err != nil {
			return "", err
		}
	}
	params.Set("apikey", i.Config.ApiKey)
	reqURL := ""
	httpClient := &http.Client{}
	resp := &http.Response{}

	switch method {
	case http.MethodPost:
		reqURL = i.Config.ApiUrl + "/api/adapter/" + service + "/" + serviceName + "/method/json"
		req, err := http.NewRequest(method, reqURL, bytes.NewBufferString(params.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err = httpClient.Do(req)
		if err != nil {
			return "", err
		}
	case http.MethodGet:
		reqURL = i.Config.ApiUrl + "/api/adapter/apikey/" + i.Config.ApiKey + "/" + service + "/" + serviceName + "/method/json"
		req, err := http.NewRequest(method, reqURL, nil)
		if err != nil {
			return "", err
		}
		req.URL.RawQuery = params.Encode()
		resp, err = httpClient.Do(req)
		if err != nil {
			return "", err
		}
	default:
		log.Errorf("method unsupported[%s]", method)
		return "", http.ErrNotSupported
	}

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(byteBody), nil
}

// Webservice queries a given webservice with all params supplied
// Returns err != nil if query fails
func (i *InfoCMDB) Webservice(ws string, params url.Values) (string, error) {
	log.Debugf("Webservice: %s, Params: %v", ws, params)
	r, err := i.Post("query", ws, params)
	if err != nil {
		return "", err
	}
	return r, nil
}
