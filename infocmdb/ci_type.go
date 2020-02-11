package infocmdb

import (
	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
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

	cacheKey := "GetCiRelationTypeIdByRelationTypeName_" + name
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

