package client

import (
	"fmt"
	"gopkg.in/resty.v1"
)

type Client struct {
	resty *resty.Client
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

func New(baseURL string) (c *Client) {
	c = &Client{}
	c.resty = resty.New().
		SetHostURL(baseURL)
	return
}

func (c *Client) NewRequest() *resty.Request {
	return c.resty.NewRequest()
}

func (c *Client) SetAuthToken(token string) {
	c.resty.SetAuthToken(token)
}

func (res ResponseError) Error() string {
	return fmt.Sprintf("%s\n%+v", res.Message, res.Data)
}
