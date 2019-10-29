package client

import (
	"gopkg.in/resty.v1"
)

type Client struct {
	resty *resty.Client
}

// Response is the default json return of the cmdb upon any request success or error
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// ResponseError is used for the resty SetError Function as a reference to capture the error message on failure
type ResponseError struct {
	Response
}

func NewClient(baseURL string) (c *Client) {
	c = &Client{}
	c.resty = resty.New().
		SetHostURL(baseURL)
	return
}

func (i *Client) NewRequest() *resty.Request {
	return i.resty.NewRequest()
}

func (i *Client) SetAuthToken(token string) {
	i.resty.SetAuthToken(token)
}

func (res ResponseError) Error() string {
	return res.Message
}
