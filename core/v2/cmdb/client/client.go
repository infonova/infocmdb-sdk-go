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

func (i *Client) SetAuthToken(token string) {
	i.resty.SetAuthToken(token)
}
func (i *Client) SetHostURL(token string) {
	i.resty.SetHostURL(token)
}

func (i *Client) Get(url string, out interface{}, params map[string]string) (err error) {
	resp, err := i.resty.R().
		SetQueryParams(params).
		SetResult(&out).
		SetError(ErrorReturn{}).
		Get(url)

	if err != nil {
		err = errors.New(fmt.Sprintf("failed to send GET request to %s: ", url, err.Error()))
		return
	}

	// HTTP status code >= 400
	if resp.IsError() {
		errResp := resp.Error().(*ErrorReturn)
		err = errors.New(errResp.Message)
		return
	}

	// out will be set in resty
	return
}

func (i *Client) Post(url string, out interface{}, params map[string]string) (err error) {
	resp, err := i.resty.R().
		SetQueryParams(params).
		SetResult(out).
		SetError(ErrorReturn{}).
		Post(url)

	if err != nil {
		err = errors.New(fmt.Sprintf("failed to send POST request to %s: ", url, err.Error()))
		return
	}

	// HTTP status code >= 400
	if resp.IsError() {
		errResp := resp.Error().(*ErrorReturn)
		err = errors.New(errResp.Message)
		return
	}

	// out will be set in resty
	return
}
