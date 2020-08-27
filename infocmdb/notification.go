package infocmdb

import (
	v1 "github.com/infonova/infocmdb-sdk-go/infocmdb/v1/infocmdb"
)

func (c *Client) SendNotification(name string, par v1.NotifyParams) (resp v1.NotificationResponse, err error) {
	return c.v1.SendNotification(name, par)
}
