package infocmdb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
)

type Ci struct {
	CiID               int    `json:"ci_id,string"`
	CiTypeID           int    `json:"ci_type_id,string"`
	CiType             string `json:"ci_type"`
	ProjectsAsString   string `json:"project"`
	ProjectIDsAsString string `json:"project_id"`
	Projects           []string
	ProjectIDs         []int
}

type getCi struct {
	Data []struct {
		Ci
	} `json:"data"`
}

func (c *Client) GetCi(ciID int) (r Ci, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciID),
	}

	jsonRet := getCi{}
	err = c.v2.Query("int_getCi", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Debugf("Error: %v", err)
		return
	}

	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(strconv.Itoa(ciID) + " - " + v2.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Ci
	default:
		err = utilError.FunctionError(strconv.Itoa(ciID) + " - " + v2.ErrTooManyResults.Error())
	}

	r.Projects = strings.Split(r.ProjectsAsString, ",") // not safe :-/
	projectIds := strings.Split(r.ProjectIDsAsString, ",")
	for _, projectIdString := range projectIds {
		projectId, _ := strconv.Atoi(projectIdString)
		r.ProjectIDs = append(r.ProjectIDs, projectId)
	}

	return
}

type CiIds []int

type getListOfCiIdsOfCiType struct {
	Data []struct {
		CiID int `json:"ciid,string"`
	} `json:"data"`
}

func (c *Client) GetListOfCiIdsOfCiType(ciTypeID int) (ciIds CiIds, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	if ciTypeID <= 0 {
		return nil, errors.New("CiTypeID must be integer greater 0")
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciTypeID),
	}

	ret := getListOfCiIdsOfCiType{}
	err = c.v2.Query("int_getListOfCiIdsOfCiType", &ret, params)
	if err != nil {
		log.Error("Error: ", err)
		return ciIds, err
	}

	for _, ciIdOfCiType := range ret.Data {
		ciIds = append(ciIds, ciIdOfCiType.CiID)
	}

	return
}

func (c *Client) GetListOfCiIdsOfCiTypeV2(ciTypeID int) (ciIds CiIds, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	if ciTypeID <= 0 {
		return nil, errors.New("CiTypeID must be integer greater 0")
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciTypeID),
	}

	ret := getListOfCiIdsOfCiType{}
	err = c.v2.Query("int_getListOfCiIdsOfCiType", &ret, params)
	if err != nil {
		log.Error("Error: ", err)
		return ciIds, err
	}

	for _, ciIdOfCiType := range ret.Data {
		ciIds = append(ciIds, ciIdOfCiType.CiID)
	}

	return
}

func (c *Client) GetListOfCiIdsOfCiTypeName(ciTypeName string) (ciIds CiIds, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	ciTypeId, err := c.GetCiTypeIdByCiTypeName(ciTypeName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to resolve id for ciTypeName '%s': %s", ciTypeName, err.Error()))
		return
	}

	ciIds, err = c.GetListOfCiIdsOfCiType(ciTypeId)
	return
}

func (c *Client) GetAndBindListOfCiAttributesOfCiTypeName(ciTypeName string, out interface{}) (err error) {
	ciIds, err := c.GetListOfCiIdsOfCiTypeName(ciTypeName)
	if err != nil {
		err = errors.New("failed to get \"" + ciTypeName + "\" ci ids: " + err.Error())
		return
	}

	err = c.GetAndBindListOfCiAttributes(ciIds, out)
	return
}

type getListOfCiIdsByAttributeValue struct {
	Data []struct {
		CiID int `json:"ci_id,string"`
	} `json:"data"`
}

func (c *Client) GetListOfCiIdsByAttributeValue(name string, value string, valueType v2.AttributeValueType) (ciIds CiIds, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	attrId, err := c.GetAttributeIdByAttributeName(name)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"argv1": strconv.Itoa(attrId),
		"argv2": value,
		"argv3": string(valueType),
	}

	ret := getListOfCiIdsByAttributeValue{}
	err = c.v2.Query("int_getCiIdByCiAttributeValue", &ret, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	for _, ciId := range ret.Data {
		ciIds = append(ciIds, ciId.CiID)
	}

	return
}

type getListOfCiIdsByCiRelation struct {
	Data []struct {
		CiId int `json:"ci_id,string"`
	} `json:"data"`
}

func (c *Client) GetListOfCiIdsByCiRelation(ciId int, ciRelationTypeName string, direction v2.CiRelationDirection) (r CiIds, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	ciRelationTypeId, err := c.GetCiRelationTypeIdByRelationTypeName(ciRelationTypeName)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	var webservice string
	params := map[string]string{
		"argv1": strconv.Itoa(ciId),
		"argv2": strconv.Itoa(ciRelationTypeId),
	}

	switch direction {
	case v2.CI_RELATION_DIRECTION_ALL:
		webservice = "int_getListOfCiIdsByCiRelation_directionList"
		params["argv3"] = "0,1,2,3,4"
	case v2.CI_RELATION_DIRECTION_DIRECTED_FROM:
		webservice = "int_getListOfCiIdsByCiRelation_directedFrom"
	case v2.CI_RELATION_DIRECTION_DIRECTED_TO:
		webservice = "int_getListOfCiIdsByCiRelation_directedTo"
	case v2.CI_RELATION_DIRECTION_BIDIRECTIONAL:
		webservice = "int_getListOfCiIdsByCiRelation_directionList"
		params["argv3"] = "3"
	case v2.CI_RELATION_DIRECTION_OMNIDIRECTIONAL:
		webservice = "int_getListOfCiIdsByCiRelation_directionList"
		params["argv3"] = "0,4"
	}

	jsonRet := getListOfCiIdsByCiRelation{}
	err = c.v2.Query(webservice, &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	for _, row := range jsonRet.Data {
		r = append(r, row.CiId)
	}

	return
}

type CreateCi struct {
	ID        int    `json:"id,string"`
	CiTypeID  int    `json:"ci_type_id,string"`
	Icon      string `json:"icon"`
	HistoryID int    `json:"history_id,string"`
	ValidFrom string `json:"valid_from"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type createCiResponse struct {
	Data []struct {
		CreateCi
	} `json:"data"`
}

func (c *Client) CreateCi(ciTypeID int, icon string, historyID int) (r CreateCi, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciTypeID),
		"argv2": icon,
		"argv3": strconv.Itoa(historyID),
	}

	jsonRet := createCiResponse{}
	err = c.v2.Query("int_createCi", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return r, err
	}

	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(strconv.Itoa(ciTypeID) + " - " + v2.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].CreateCi
	default:
		err = utilError.FunctionError(strconv.Itoa(ciTypeID) + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

func (c *Client) GetCiIdByAttributeValue(name string, value string, valueType v2.AttributeValueType) (ciId int, err error) {
	ciIds, err := c.GetListOfCiIdsByAttributeValue(name, value, valueType)
	if err != nil {
	    return
	}

	switch len(ciIds) {
	case 0:
		err = utilError.FunctionError(name + " - " + v2.ErrNoResult.Error())
	case 1:
		ciId = ciIds[0]
	default:
		err = utilError.FunctionError(name + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

func (c *Client) GetAndBindCiByAttributeValue(name string, value string, valueType v2.AttributeValueType, out interface{}) (err error) {
	ciId, err := c.GetCiIdByAttributeValue(name, value, valueType)
	if err != nil {
	    return
	}

	err = c.GetAndBindCiAttributes(ciId, out)
	if err != nil {
	    return
	}

	return
}
