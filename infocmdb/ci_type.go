package infocmdb

import (
	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"

	"strconv"

	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type getCiTypeIdByCiTypeName struct {
	Data []responseId `json:"data"`
}

func (c *Client) GetCiTypeIdByCiTypeName(name string) (r int, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetCiTypeIdByCiTypeName_" + name
	cached, found := c.v1.Cache.Get(cacheKey)
	if found {
		return cached.(int), nil
	}

	params := map[string]string{
		"argv1": name,
	}

	response := getCiTypeIdByCiTypeName{}
	err = c.v2.Query("int_getCiTypeIdByCiTypeName", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(name + " - " + v2.ErrNoResult.Error())
	case 1:
		r = response.Data[0].Id
		c.v1.Cache.Set(cacheKey, r, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

type getCiTypeName struct {
	Data []struct {
		Name string `json:"name"`
	} `json:"data"`
}

func (c *Client) GetCiTypeName(ciId int) (ciTypeName string, err error) {
	ciIdString := strconv.Itoa(ciId)

	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetCiTypeName_" + ciIdString
	cached, found := c.v1.Cache.Get(cacheKey)
	if found {
		return cached.(string), nil
	}

	params := map[string]string{
		"argv1": ciIdString,
		"argv2": "name",
	}
	response := getCiTypeName{}
	err = c.v2.Query("int_getCiTypeOfCi", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(ciIdString + " - " + v2.ErrNoResult.Error())
	case 1:
		ciTypeName = response.Data[0].Name
		c.v1.Cache.Set(cacheKey, ciTypeName, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(ciIdString + " - " + v2.ErrTooManyResults.Error())
	}

	return
}
