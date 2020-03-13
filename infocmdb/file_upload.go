package infocmdb

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	"github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb/client"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
)

type uploadFileResponse struct {
	Data string `json:"data"`
}

// Supported file data types are `string`, `[]byte` and `io.Reader`.
// The returned uploadId can be used for attachment attributes in the `UpdateCiAttribute` function.
func (c *Client) UploadFile(file interface{}) (uploadId string, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	var response uploadFileResponse
	var respErr client.ResponseError

	resp, err := c.v2.Client.Execute(resty.MethodPost, "/apiV2/fileupload",
		func(request *resty.Request) *resty.Request {
			return request.
				SetBody(file).
				SetResult(&response).
				SetError(&respErr)
		})

	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	if resp != nil && resp.IsError() {
		err = respErr
		return
	}

	uploadId = response.Data
	return
}
