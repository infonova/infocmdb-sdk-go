package client

import (
	"errors"
	"github.com/go-resty/resty"
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
	Message string `json:"message"`
	Success bool `json:"success"`
	Data interface{} `json:"data,omitempty"`
}
func (i *Client) SetAuthToken(token string) {
	i.resty.SetAuthToken(token)
}
func (i *Client) SetHostURL(token string) {
	i.resty.SetHostURL(token)
}

func (i *Client) Get(url string, out interface {}, params map[string]string) (err error) {
	resp, err := i.resty.R().
		SetQueryParams(params).
		SetResult(&out).
		SetError(ErrorReturn{}).
		Get(url)

	if err != nil {
		return
	}
	errResp := resp.Error().(*ErrorReturn)
	if !errResp.Success {
		return errors.New(errResp.Message)
	}
	return
}

func (i *Client) Post(url string, out interface {}, params map[string]string) (err error) {
	resp, err := i.resty.R().
		SetQueryParams(params).
		SetResult(&out).
		SetError(ErrorReturn{}).
		Post(url)

	if err != nil {
		return
	}
	errResp := resp.Error().(*ErrorReturn)
	if !errResp.Success {
		return errors.New(errResp.Message)
	}
	return
}