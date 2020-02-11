package infocmdb

import (
	"strconv"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type getProjectIdByProjectName struct {
	Data []responseId `json:"data"`
}

func (c *Client) GetProjectIdByProjectName(name string) (projectID int, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetProjectIdByProjectName_" + name
	cached, found := c.v2.Cache.Get(cacheKey)
	if found {
		return cached.(int), nil
	}

	params := map[string]string{
		"argv1": name,
	}

	response := getProjectIdByProjectName{}
	err = c.v2.Query("int_getProjectIdByProjectName", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(name + " - " + v2.ErrNoResult.Error())
	case 1:
		projectID = response.Data[0].Id
		c.v2.Cache.Set(cacheKey, projectID, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

type addCiProjectMappingResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    []struct{} `json:"data"`
}

func (c *Client) AddCiProjectMapping(ciID int, projectID int, historyID int) (err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciID),
		"argv2": strconv.Itoa(projectID),
		"argv3": strconv.Itoa(historyID),
	}

	jsonRet := addCiProjectMappingResponse{}
	err = c.v2.Query("int_addCiProjectMapping", &jsonRet, params)
	if err != nil {
		log.Error("Error: ", err)
	}

	return
}
