package client

import (
	"errors"
	"fmt"
	"gopkg.in/resty.v1"
)

type Client struct {
	resty *resty.Client
}

func NewClient(baseURL string) (c *Client) {
	c = &Client{}
	c.resty = resty.New().
		SetHostURL(baseURL)
	return
}

type ErrorReturn struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

const MethodPostForm = "POST_FORM"

func (i *Client) SetAuthToken(token string) {
	i.resty.SetAuthToken(token)
}
func (i *Client) SetHostURL(token string) {
	i.resty.SetHostURL(token)
}

func (i *Client) Execute(method string, url string, out interface{}, params map[string]string) (err error) {
	rest := i.resty.R().
		SetResult(&out).
		SetError(ErrorReturn{})

	switch method {
	case resty.MethodGet:
		rest.SetQueryParams(params)
	case MethodPostForm:
		rest.SetFormData(params)
		method = resty.MethodPost
	default:
		rest.SetBody(params) // json encoded body
	}

	resp, err := rest.Execute(method, url)

	if err != nil {
		err = errors.New(fmt.Sprintf("failed to send %s request to %s: %s", method, url, err.Error()))
		return
	}

	// HTTP status code >= 400
	if resp.IsError() {
		errResp := resp.Error().(*ErrorReturn)
		err = errors.New(errResp.Message)
		return
	}

	// "out" will be set in resty
	return
}

func (i *Client) Get(url string, out interface{}, params map[string]string) (err error) {
	return i.Execute(resty.MethodGet, url, out, params)
}

func (i *Client) Post(url string, out interface{}, params map[string]string) (err error) {
	return i.Execute(resty.MethodPost, url, out, params)
}

func (i *Client) PostForm(url string, out interface{}, params map[string]string) (err error) {
	return i.Execute(MethodPostForm, url, out, params)
}

func (i *Client) Put(url string, out interface{}, params map[string]string) (err error) {
	return i.Execute(resty.MethodPut, url, out, params)
}

func (i *Client) Delete(url string, out interface{}, params map[string]string) (err error) {
	return i.Execute(resty.MethodDelete, url, out, params)
}

func (i *Client) Patch(url string, out interface{}, params map[string]string) (err error) {
	return i.Execute(resty.MethodPatch, url, out, params)
}
