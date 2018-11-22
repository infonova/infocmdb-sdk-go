package infocmdbGoLib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type CmdbWebClient struct {
	apikey string
	url    string
	//client http.Client
}

func NewCmdbWebClient() CmdbWebClient {
	return CmdbWebClient{}
}

type jsonLoginReturn struct {
	Status string `json:"status"`
	ApiKey string `json:"apikey"`
}

func (c *CmdbWebClient) Login(url string, username string, password string) error {
	log.Printf("Opening new Webclient connection. (Url: %s, Username: %s)", url, username)

	c.url = url

	reqURL := fmt.Sprintf("%s/api/login/username/%s/password/%s/timeout/21600/method/json", url, username, password)

	resp, err := http.Get(reqURL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response.", err)
		return err
	}

	var loginResult jsonLoginReturn
	err = json.Unmarshal(byteBody, &loginResult)
	if err != nil {
		log.Println("Error unmarshaling Json.", err)
		return err
	}

	if loginResult.Status != "OK" {
		log.Printf("Login Status not ok. Status: \"%s\"\n", loginResult.Status)
		return errors.New("Login Status not ok.")
	}

	c.apikey = loginResult.ApiKey
	return nil
}

func (c *CmdbWebClient) LoginWithApiKey(url string, apikey string) error {
	log.Printf("Opening new Webclient connection using Apikey. (Url: %s, ApiKey: %s)", url, apikey)
	c.url = url
	c.apikey = apikey
	return nil
}

// CallGet the api with a given method and parameters as a GetRequest
func (c *CmdbWebClient) CallGet(service string, serviceName string, params url.Values) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET",
		c.url+"/api/adapter/apikey/"+c.apikey+"/"+service+"/"+serviceName+"/method/json", nil)
	if err != nil {
		log.Print(err)
		return "", err
	}
	if params.Get("apikey") == "" {
		params.Set("apikey", c.apikey)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response.", err)
		return "", err
	}

	return string(byteBody), nil
}

// CallPost the api with a given method and parameters as a GetRequest
func (c *CmdbWebClient) CallPost(service string, serviceName string, postData url.Values) (string, error) {
	//client := &http.Client{	}
	reqURL := c.url + "/api/adapter/" + service + "/" + serviceName + "/method/json"
	if postData.Get("apikey") == "" {
		postData.Set("apikey", c.apikey)
	}
	log.Println(postData)
	resp, err := http.PostForm(reqURL, postData)
	//req, err := http.NewRequest("POST", , nil)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response.", err)
		return "", err
	}

	return string(byteBody), nil
}
