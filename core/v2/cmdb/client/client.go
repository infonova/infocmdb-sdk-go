package client

import (
	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb/models"
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

func (i *Client) NewRequest() *resty.Request {
	r := i.resty.R().
		SetError(models.ErrorReturn{})

	return r
}

func (i *Client) SetAuthToken(token string) {
	i.resty.SetAuthToken(token)
}
