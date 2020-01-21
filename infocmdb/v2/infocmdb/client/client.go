package client

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"strconv"
	"strings"
)

type LoginParams struct {
	Username string
	Password string
	Lifetime int
}

type Client struct {
	resty       *resty.Client
	loginParams LoginParams
}

// Response is the default json return of the cmdb upon any request success or error
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseError is used for the resty SetError Function as a reference to capture the error message on failure
type ResponseError struct {
	Response
}

func (res ResponseError) Error() string {
	return fmt.Sprintf("%s\n%+v", res.Message, res.Data)
}

func New(baseURL string) (c *Client) {
	c = &Client{}
	c.resty = resty.New().
		SetHostURL(baseURL)
	return
}

type loginTokenReturn struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (c *Client) Login(loginParams LoginParams) (token string, err error) {
	c.loginParams = loginParams

	if loginParams.Username == "" || loginParams.Password == "" {
		return "", errors.New("must provide credentials")
	}

	var loginResult loginTokenReturn
	params := map[string]string{
		"username": loginParams.Username,
		"password": loginParams.Password,
		"lifetime": strconv.Itoa(loginParams.Lifetime),
	}

	var errResp ResponseError
	resp, err := c.resty.
		NewRequest().
		SetError(&errResp).
		SetResult(&loginResult).
		SetFormData(params).
		Post("/apiV2/auth/token")

	if err != nil {
		return "", err
	}

	if resp != nil && resp.IsError() {
		return "", errResp
	}

	token = loginResult.Data.Token

	if token == "" {
		return "", errors.New("login status not ok")
	}

	c.resty.SetAuthToken(token)
	return
}

type PrepareRequestFunc func(request *resty.Request) *resty.Request

// Executes a request, automatically resolving timed out API token problems and retrying.
func (c *Client) Execute(method, url string, prepareRequestFunc PrepareRequestFunc) (resp *resty.Response, err error) {
	req := prepareRequestFunc(c.resty.NewRequest())
	resp, err = req.Execute(method, url)

	if err != nil {
		return
	}

	if resp.IsError() &&
		resp.StatusCode() == 403 &&
		strings.Contains(resp.String(), "Not authenticated") {

		log.Debug("Request failed due to authentication error, logging in again and retrying...")
		token, err := c.Login(c.loginParams)
		if err != nil {
			return nil, errors.New("re-login after authentication error failed: " + err.Error())
		}

		req.SetAuthToken(token)
		return req.Execute(method, url)
	}

	return
}
