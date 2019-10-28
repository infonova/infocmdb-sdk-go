package client

import (
	"gopkg.in/resty.v1"
)

type Client struct {
	resty *resty.Client
}

type ResponseStatus struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
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

func (res ResponseStatus) Error() string {
	return res.Message
}
