package infocmdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type NotifyParams struct {
	From             string
	FromName         string
	Recipients       []string
	RecipientsCC     []string
	RecipientsBCC    []string
	Subject          string
	AttachmentsPaths []string
	OtherParams      map[string]string
}

type SendNotificationData struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type sendNotificationResp struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []SendNotificationData `json:"data"`
}

type NotificationResponse struct {
	SentTo []SendNotificationData
}

func (i *Cmdb) SendNotification(notifyName string, params NotifyParams) (resp NotificationResponse, err error) {

	err = i.Login()
	if err != nil {
		return resp, err
	}

	httpClient := &http.Client{}
	reqParams := url.Values{}

	reqParams.Add("apikey", i.Config.ApiKey)

	for bodyKey, bodyVal := range params.OtherParams {
		reqParams.Add(bodyKey, bodyVal)
	}

	reqUrl := i.Config.ApiUrl + "/api/notification/notify/" + notifyName + "/method/json"
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBufferString(reqParams.Encode()))

	if err != nil {
		return resp, errors.New("failed to create a new request. Error: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if params.From != "" {
		req.Header.Add("From", params.From)
	}

	if params.FromName != "" {
		req.Header.Add("FromName", params.FromName)
	}

	allRecipients := strings.Join(params.Recipients, ";")

	if allRecipients != "" {
		req.Header.Add("Recipients", allRecipients)
	}

	recipientsCC := strings.Join(params.RecipientsCC, ";")

	if recipientsCC != "" {
		req.Header["RecipientsCC"] = []string{recipientsCC}
	}

	recipientsBCC := strings.Join(params.RecipientsBCC, ";")

	if recipientsBCC != "" {
		req.Header.Add("RecipientsBCC", recipientsBCC)
	}

	allAttachments := strings.Join(params.AttachmentsPaths, ";")

	if allAttachments != "" {
		req.Header.Add("Attachments", allAttachments)
	}

	if params.Subject != "" {
		req.Header.Add("Subject", params.Subject)
	}

	response, err := httpClient.Do(req)

	if err != nil {
		return resp, errors.New("failed to make a request. Error: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return resp, errors.New("failed to read response body: " + err.Error())
	}

	var responseStruct sendNotificationResp
	err = json.Unmarshal(bodyBytes, &responseStruct)
	if err != nil {
		return resp, errors.New("failed to parse response body. Error: " + err.Error() + "Body: " + string(bodyBytes))
	}

	switch responseStruct.Status {
	case "OK":

		resp := NotificationResponse{
			SentTo: responseStruct.Data,
		}
		return resp, nil

	case "error":
		return resp, errors.New(responseStruct.Message)
	}

	return
}
