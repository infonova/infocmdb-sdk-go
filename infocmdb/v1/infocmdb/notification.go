package infocmdb

import (
	log "github.com/sirupsen/logrus"
)

type notificationParams struct{
	Recipients [] string
	NotificationAttributes map[string]string
}

type notificationResponse struct {
	Data string `json:"data"`
}


func (c *Cmdb) sendNotification (notificationName string, params notificationParams) (r notificationResponse, err error) {

	err = c.Login()
	if err != nil {
		log.Error(err)
		return
	}



	return
}