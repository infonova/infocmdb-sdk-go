package cmdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var (
	ErrLoginFailed = errors.New("login status not ok")
)

type ResultLogin struct {
	Status string `json:"status"`
	ApiKey string `json:"apikey"`
}

func (i *Cmdb) LoginWithUserPass(url string, username string, password string) error {

	log.Debugf("Opening new WebClient connection. (Url: %s, Username: %s)", url, username)

	reqURL := fmt.Sprintf("%s/api/login/username/%s/password/%s/timeout/21600/method/json", url, username, password)

	resp, err := http.Get(reqURL)
	if err != nil {
		errMsg := strings.Replace(err.Error(), password, "*******", -1)
		return errors.New(errMsg)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()

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

func (i *Cmdb) LoginWithApiKey(url string, apikey string) error {
	log.Debugf("Opening new WebClient connection using ApiKey. (Url: %s, ApiKey: %s)", url, apikey)
	i.Config.ApiUrl = url
	i.Config.ApiKey = apikey
	return nil
}

// validWebserviceMethod verifies that only supported (GET, POST) method calls get send to the CMDB
func validWebserviceMethod(method string) (err error) {
	switch method {
	case http.MethodGet:
		return
	case http.MethodPost:
		return
	}

	return http.ErrNotSupported
}

func (i *Cmdb) CallWebservice(method string, service string, serviceName string, params url.Values, variable interface{}) (err error) {
	if err = validWebserviceMethod(method); err != nil {
		return err
	}

	if i.Config.ApiKey == "" {
		err = i.Login()
		if err != nil {
			return err
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
			return err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err = httpClient.Do(req)
		if err != nil {
			return err
		}
	case http.MethodGet:
		reqURL = i.Config.ApiUrl + "/api/adapter/apikey/" + i.Config.ApiKey + "/" + service + "/" + serviceName + "/method/json"
		req, err := http.NewRequest(method, reqURL, nil)
		if err != nil {
			return err
		}
		req.URL.RawQuery = params.Encode()
		resp, err = httpClient.Do(req)
		if err != nil {
			return err
		}
	default:
		log.Errorf("method unsupported[%s]", method)
		return http.ErrNotSupported
	}

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = checkResponseStatusMessage(byteBody)
	if err != nil {
		return err
	}

	rv := reflect.Indirect(reflect.ValueOf(variable))
	if rv.Kind() == reflect.String {
		rv.SetString(string(byteBody))
		return nil
	}

	err = json.Unmarshal([]byte(byteBody), &variable)
	if err != nil {
		err = errors.New(err.Error() + ": " + string(byteBody))
		return err
	}
	return nil
}

type ResponseStatus struct {
	Status string `json:"status"`
}

// checkResponseStatusMessage will verify that the cmdb returend "status": "OK"
func checkResponseStatusMessage(byteBody []byte) (err error) {
	responseStatus := ResponseStatus{}
	err = json.Unmarshal([]byte(byteBody), &responseStatus)
	if err != nil {
		log.Error("error checking StatusMessage: ", err)
		log.Error(string(byteBody))
		return err
	}

	if strings.Compare(responseStatus.Status, "OK") != 0 {
		err = ErrWebserviceResponseNotOk
	}
	return
}

// Webservice queries a given webservice with all params supplied
// Returns err != nil if query fails
func (i *Cmdb) Webservice(ws string, params url.Values) (r string, err error) {
	log.Debugf("Webservice: %s, Params: %v", ws, params)
	err = i.CallWebservice(http.MethodPost, "query", ws, params, &r)
	if err != nil {
		return "", err
	}
	return r, nil
}
