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

type cmdbWebClient struct {
	apikey string
	url    string
	//client http.Client
}

func NewCmdbWebClient() cmdbWebClient {
	return cmdbWebClient{}
}

type ResultLogin struct {
	Status string `json:"status"`
	ApiKey string `json:"apikey"`
}

func (c *cmdbWebClient) Login(url string, username string, password string) error {
	log.Debug("Opening new Webclient connection. (Url: %s, Username: %s)", url, username)

	c.url = url

	reqURL := fmt.Sprintf("%s/api/login/username/%s/password/%s/timeout/21600/method/json", url, username, password)

	resp, err := http.Get(reqURL)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response.", err)
		return err
	}

	var loginResult ResultLogin
	err = json.Unmarshal(byteBody, &loginResult)
	if err != nil {
		log.Error("Error unmarshaling Json.", err)
		return err
	}

	if loginResult.Status != "OK" {
		log.Error("Login Status not ok. Status: \"%s\"\n", loginResult.Status)
		return errors.New("Login Status not ok.")
	}

	c.apikey = loginResult.ApiKey
	return nil
}

func (c *cmdbWebClient) LoginWithApiKey(url string, apikey string) error {
	log.Info("Opening new Webclient connection using Apikey. (Url: %s, ApiKey: %s)", url, apikey)
	c.url = url
	c.apikey = apikey
	return nil
}

// Get the api with a given method and parameters as a GetRequest
func (c *cmdbWebClient) Get(service string, serviceName string, params url.Values) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET",
		c.url+"/api/adapter/apikey/"+c.apikey+"/"+service+"/"+serviceName+"/method/json", nil)

	if err != nil {
		log.Error(err)
		return "", err
	}
	if params.Get("apikey") == "" {
		params.Set("apikey", c.apikey)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response.", err)
		return "", err
	}

	return string(byteBody), nil
}

// Post the api with a given method and parameters as a GetRequest
func (c *cmdbWebClient) Post(service string, serviceName string, postData url.Values) (string, error) {
	log.Debug("service: %s name: %s post: %v\n", service, serviceName, postData)
	//client := &http.Client{	}
	reqURL := c.url + "/api/adapter/" + service + "/" + serviceName + "/method/json"
	if postData.Get("apikey") == "" {
		postData.Set("apikey", c.apikey)
	}
	log.Debug(postData)
	resp, err := http.PostForm(reqURL, postData)
	//req, err := http.NewRequest("POST", , nil)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response.", err)
		return "", err
	}

	return string(byteBody), nil
}
