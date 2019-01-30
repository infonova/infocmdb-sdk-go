package infocmdblibrary

import (
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

type CmdbWebClient struct {
	ApiKey string
	Url    string
	//client http.Client
}

func NewCmdbWebClient() CmdbWebClient {
	return CmdbWebClient{}
}

type ResultLogin struct {
	Status string `json:"status"`
	ApiKey string `json:"apikey"`
}

func (c *CmdbWebClient) Login(url string, username string, password string) error {
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

	c.ApiKey = loginResult.ApiKey
	return nil
}

func (c *CmdbWebClient) LoginWithApiKey(url string, apikey string) error {
	log.Debugf("Opening new Webclient connection using Apikey. (Url: %s, ApiKey: %s)", url, apikey)
	c.Url = url
	c.ApiKey = apikey
	return nil
}

// Get the api with a given method and parameters as a GetRequest
func (c *CmdbWebClient) Get(service string, serviceName string, params url.Values) (string, error) {
	params.Set("apikey", c.ApiKey)
	log.Debugf("service: %s name: %s post: %v", service, serviceName, params)
	log.Debugf("Data: %v", params)

	reqURL := c.Url + "/api/adapter/apikey/" + c.ApiKey + "/" + service + "/" + serviceName + "/method/json"

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return string(byteBody), nil
}

// Post the api with a given method and parameters as a GetRequest
func (c *CmdbWebClient) Post(service string, serviceName string, params url.Values) (string, error) {
	params.Set("apikey", c.ApiKey)
	log.Debugf("service: %s name: %s post: %v", service, serviceName, params)
	log.Debugf("Data: %v", params)

	reqURL := c.Url + "/api/adapter/" + service + "/" + serviceName + "/method/json"

	resp, err := http.PostForm(reqURL, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(byteBody), nil
}
