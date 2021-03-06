package infocmdb

import (
	"regexp"
	"strconv"
	"strings"

	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
)

type createCiRelation struct {
}

func (c *Client) CreateCiRelation(ciId1 int, ciId2 int, ciRelationTypeName string, direction v2.CiRelationDirection) (err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	directionId, err := direction.GetId()
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	counter, err := c.GetCiRelationCount(ciId1, ciId2, ciRelationTypeName)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	if counter == 0 {
		var ciRelationTypeId int
		ciRelationTypeId, err = c.GetCiRelationTypeIdByRelationTypeName(ciRelationTypeName)
		if err != nil {
			err = utilError.FunctionError(err.Error())
			return
		}

		params := map[string]string{
			"argv1": strconv.Itoa(ciId1),
			"argv2": strconv.Itoa(ciId2),
			"argv3": strconv.Itoa(ciRelationTypeId),
			"argv4": strconv.Itoa(directionId),
		}

		jsonRet := createCiRelation{}
		err = c.v2.Query("int_createCiRelation", &jsonRet, params)
		if err != nil {
			err = utilError.FunctionError(err.Error())
			log.Error("Error: ", err)
			return
		}
	}

	return
}

type deleteCiRelation struct {
}

func (c *Client) DeleteCiRelation(ciId1 int, ciId2 int, ciRelationTypeName string) (err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	ciRelationTypeId, err := c.GetCiRelationTypeIdByRelationTypeName(ciRelationTypeName)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciId1),
		"argv2": strconv.Itoa(ciId2),
		"argv3": strconv.Itoa(ciRelationTypeId),
	}

	jsonRet := deleteCiRelation{}
	err = c.v2.Query("int_deleteCiRelation", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	return
}

func (c *Client) AttributeBasedRelation(sourceCiId int, attributeName string, ciRelationTypeName string, triggerType string, swapCiColumns bool) (relationCisAdded []int, relationCisRemoved []int, err error) {
	var currentCiValues []string
	if triggerType != "ci_attribute_delete" {
		var value string
		value, _, err = c.GetCiAttributeValueCi(sourceCiId, attributeName)
		if err != nil {
			if !strings.Contains(err.Error(), v2.ErrNoResult.Error()) {
				return
			}
		} else if value != "" {
			currentCiValues = regexp.MustCompile(`,\s?`).Split(value, -1)
		}
	}

	destinationCiIds := make([]int, len(currentCiValues))
	for i, s := range currentCiValues {
		destinationCiIds[i], _ = strconv.Atoi(strings.TrimSpace(s))
	}

	return c.CiBasedRelation(sourceCiId, destinationCiIds, ciRelationTypeName, triggerType, swapCiColumns)
}

//noinspection GoUnusedParameter preserved to remain backwards compatible
func (c *Client) CiBasedRelation(srcCiId int, destCiId []int, ciRelationTypeName string, triggerType string, swapCiColumns bool) (relationCisAdded []int, relationCisRemoved []int, err error) {
	currentCiRelations, err := c.GetListOfCiIdsByCiRelation(srcCiId, ciRelationTypeName, v2.CI_RELATION_DIRECTION_ALL)
	if err != nil {
		return
	}

	// add relations
	for _, valueCiId := range destCiId {
		add := true
		for _, relationCiId := range currentCiRelations {
			if valueCiId == relationCiId {
				add = false
				break
			}
		}

		if add {
			if swapCiColumns {
				err = c.CreateCiRelation(valueCiId, srcCiId, ciRelationTypeName, v2.CI_RELATION_DIRECTION_OMNIDIRECTIONAL)
			} else {
				err = c.CreateCiRelation(srcCiId, valueCiId, ciRelationTypeName, v2.CI_RELATION_DIRECTION_OMNIDIRECTIONAL)
			}
			if err != nil {
				return
			}

			relationCisAdded = append(relationCisAdded, valueCiId)
		}
	}

	// remove relations
	for _, relationCiId := range currentCiRelations {
		remove := true
		for _, valueCiId := range destCiId {
			if valueCiId == relationCiId {
				remove = false
				break
			}
		}

		if remove {
			err = c.DeleteCiRelation(srcCiId, relationCiId, ciRelationTypeName)
			if err != nil {
				return
			}

			relationCisRemoved = append(relationCisRemoved, relationCiId)
		}
	}

	return
}

type getCiRelationCount struct {
	Data []struct {
		Count int `json:"c,string"`
	} `json:"data"`
}

func (c *Client) GetCiRelationCount(ciId1 int, ciId2 int, ciRelationTypeName string) (r int, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	ciRelationTypeId, err := c.GetCiRelationTypeIdByRelationTypeName(ciRelationTypeName)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciId1),
		"argv2": strconv.Itoa(ciId2),
		"argv3": strconv.Itoa(ciRelationTypeId),
	}

	jsonRet := getCiRelationCount{}
	err = c.v2.Query("int_getCiRelationCount", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return r, err
	}

	errPrefix := strconv.Itoa(ciId1) + ", " + strconv.Itoa(ciId2) + ", " + ciRelationTypeName + "(" + strconv.Itoa(ciRelationTypeId) + ")" + " - "
	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(errPrefix + v2.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Count
	default:
		err = utilError.FunctionError(errPrefix + v2.ErrTooManyResults.Error())
	}

	return

}

type getCiRelationTypeIdByRelationTypeName struct {
	Data []responseId `json:"data"`
}

func (c *Client) GetCiRelationTypeIdByRelationTypeName(name string) (r int, err error) {
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

	jsonRet := getCiRelationTypeIdByRelationTypeName{}
	err = c.v2.Query("int_getCiRelationTypeIdByRelationTypeName", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(name + " - " + v2.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Id
		c.v1.Cache.Set(cacheKey, r, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

type Relation struct {
	Id        int
	CiId1     int
	CiId2     int
	Direction v2.CiRelationDirection
}

type getCiRelationsByName struct {
	Data []struct {
		Id        IntWrapper `json:"id"`
		CiId1     IntWrapper `json:"ci_id_1"`
		CiId2     IntWrapper `json:"ci_id_2"`
		Direction IntWrapper `json:"direction"`
	} `json:"data"`
}

func (c *Client) GetListOfRelationsByName(name string) (relations []Relation, err error) {
	relations = []Relation{}

	if err = c.v2.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": name,
	}

	jsonRet := getCiRelationsByName{}
	err = c.v2.Query("int_getCiRelationsByName", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	for _, relation := range jsonRet.Data {
		direction, err := v2.NewCiRelationDirection(int(relation.Direction))
		if err != nil {
		    return relations, err
		}

		relations = append(relations, Relation{
			Id:        int(relation.Id),
			CiId1:     int(relation.CiId1),
			CiId2:     int(relation.CiId2),
			Direction: direction,
		})
	}

	return
}
