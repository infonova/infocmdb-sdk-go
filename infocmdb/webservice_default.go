package infocmdb

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	v1 "github.com/infonova/infocmdb-sdk-go/infocmdb/v1/cmdb"
	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/cmdb"
	clientV2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/cmdb/client"

	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// QueryWebservices allows you to call a generic webservice(arg1: ws) with the providing params
// Return: json string
func (c *Client) QueryWebservice(ws string, params map[string]string) (resp string, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	resp, err = c.v2.QueryRaw(ws, params)
	if err != nil {
		log.Error("Error: ", err)
	}

	return
}

// Query allows you to call a generic webservice(arg1: ws) with the providing params and a reference
// to a result. It will take the built in resty function to deserialize the result
// Return: error
func (c *Client) Query(ws string, out interface{}, params map[string]string) (err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	if err = c.v2.Query(ws, out, params); err != nil {
		log.Error("Error: ", err)
	}

	return
}

type CiIds []int

type GetListOfCiIdsOfCiType struct {
	Status string `json:"status"`
	Data   []struct {
		CiID int `json:"ciid,string"`
	} `json:"data"`
}

type ResponseId struct {
	Id int `json:"id,string"`
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

	ret := GetListOfCiIdsOfCiType{}
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

	ret := GetListOfCiIdsOfCiType{}
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

type GetListOfCiIdsByAttributeValue struct {
	Data []struct {
		CiID int `json:"ci_id,string"`
	} `json:"data"`
}

func (c *Client) GetListOfCiIdsByAttributeValue(att, value string, field v2.AttributeValueType) (ciIds CiIds, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	attId, err := c.GetAttributeIdByAttributeName(att)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"argv1": strconv.Itoa(attId),
		"argv2": value,
		"argv3": string(field),
	}

	ret := GetListOfCiIdsByAttributeValue{}
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

// // Templates for others
type AddCiProjectMapping struct {
	Status string `json:"status"`
}

// AddCiProjectMapping
// int_addCiProjectMapping     add project-mapping to a ci
//
// insert into ci_project (ci_id, project_id, history_id)
// select :argv1:, :argv2:, :argv3:
// from dual
// where not exists(select id from ci_project where ci_id = :argv1: and project_id = :argv2:)

func (c *Client) AddCiProjectMapping(ciID int, projectID int, historyID int) (r AddCiProjectMapping, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
		"argv2": {strconv.Itoa(projectID)},
		"argv3": {strconv.Itoa(historyID)},
	}

	err = c.v1.CallWebservice(http.MethodPost, "query", "int_addCiProjectMapping", params, &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// // CreateAttribute
// // int_createAttribute     create an attribute
type CreateAttribute struct {
	Status string `json:"status"`
	CiID   int    `json:"id"`
}

func (c *Client) CreateAttribute(ciID int, attrID int) (r CreateAttribute, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
		"argv2": {strconv.Itoa(attrID)},
	}

	err = c.v1.CallWebservice(http.MethodPost, "query", "int_createCiAttribute", params, &r)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// CreateCi
// int_createCi    create a CI
type CreateCi struct {
	Status    string `json:"status"`
	CiTypeID  int    `json:"ci_type_id"`
	Icon      string `json:"icon"`
	ValidFrom string `json:"valid_from"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *Client) CreateCi(ciTypeID int, icon string, historyID int) (r CreateCi, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	params := url.Values{
		"argv1": {strconv.Itoa(ciTypeID)},
		"argv2": {icon},
		"argv3": {strconv.Itoa(historyID)},
	}

	err = c.v1.CallWebservice(http.MethodPost, "query", "int_createCi", params, &r)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return r, err
	}
	return
}

// GetCi
// int_getCi   Retrieve all informations about a ci
type Ci struct {
	CiID               int    `json:"ci_id,string"`
	CiTypeID           int    `json:"ci_type_id,string"`
	CiType             string `json:"ci_type"`
	ProjectsAsString   string `json:"project"`
	ProjectIDsAsString string `json:"project_id"`
	Projects           []string
	ProjectIDs         []int
}
type GetCi struct {
	Status string `json:"status"`
	Data   []struct {
		Ci
	} `json:"data"`
}

func (c *Client) GetCi(ciID int) (r Ci, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
	}

	jsonRet := GetCi{}
	err = c.v1.CallWebservice(http.MethodPost, "query", "int_getCi", params, &jsonRet)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Debugf("Error: %v", err)
		return
	}

	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(strconv.Itoa(ciID) + " - " + v1.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Ci
	default:
		err = utilError.FunctionError(strconv.Itoa(ciID) + " - " + v1.ErrTooManyResults.Error())
	}

	r.Projects = strings.Split(r.ProjectsAsString, ",") // not safe :-/
	projectIds := strings.Split(r.ProjectIDsAsString, ",")
	for _, projectIdString := range projectIds {
		projectId, _ := strconv.Atoi(projectIdString)
		r.ProjectIDs = append(r.ProjectIDs, projectId)
	}

	return
}

// GetCiAttributes
// int_getCiAttributes     get all attributes for given ci (:argv1:)

type CiAttribute struct {
	CiID                 int    `json:"ci_id,string,string"`
	CiAttributeID        int    `json:"ci_attribute_id,string"`
	AttributeID          int    `json:"attribute_id,string"`
	AttributeName        string `json:"attribute_name"`
	AttributeDescription string `json:"attribute_description"`
	AttributeType        string `json:"attribute_type"`
	Value                string `json:"value"`
	ModifiedAt           string `json:"modified_at"`
}

type GetCiAttributes struct {
	Status string        `json:"status"`
	Data   []CiAttribute `json:"data"`
}

func (c *Client) GetCiAttributes(ciID int) (r []CiAttribute, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
	}

	jsonRet := GetCiAttributes{}
	err = c.v1.CallWebservice(http.MethodPost, "query", "int_getCiAttributes", params, &jsonRet)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
	}
	r = jsonRet.Data
	return
}

// CreateCiAttribute
// int_createCiAttribute   creates a ci_attribute-row
type CreateCiAttribute struct {
	Status      string `json:"status"`
	CiID        int    `json:"ci_id"`
	AttributeID int    `json:"attribute_id"`
	HistoryID   int    `json:"history_id"`
}

func (c *Client) CreateCiAttribute() (r CreateCiAttribute, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	/*
		params := url.Values{
			// "argv1": {strconv.Itoa(%PARAM1%)},
			// "argv2": {strconv.Itoa(%PARAM2%)},
			// "argv3": {strconv.Itoa(%PARAM3%)},
			// "argv4": {strconv.Itoa(%PARAM4%)},
		}

		ret, err := c.v1.CallWebservice(http.MethodPost,"query", "int_createCiAttribute", params)
		if err != nil {
			log.Error("Error: ", err)
			return r, err
		}

		err = json.Unmarshal([]byte(ret), &r)
		if err != nil {
			log.Error("Error: ", err)
			return r, err
		}

		return
	*/
}

// CreateCiRelation
// int_createCiRelation    inserts a relation: argv1 = ci_id_1 argv2 = ci_id_2 argv3 = ci_relation_type_id argv4 = direction
type CreateCiRelation struct {
	Status string `json:"status"`
}

func (c *Client) CreateCiRelation(ciId1 int, ciId2 int, ciRelationTypeName string, direction v1.CiRelationDirection) (err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	var directionId int
	switch direction {
	case v1.CI_RELATION_DIRECTION_DIRECTED_FROM:
		directionId = 1
	case v1.CI_RELATION_DIRECTION_DIRECTED_TO:
		directionId = 2
	case v1.CI_RELATION_DIRECTION_BIDIRECTIONAL:
		directionId = 3
	case v1.CI_RELATION_DIRECTION_OMNIDIRECTIONAL:
		directionId = 4
	default:
		err = errors.New(fmt.Sprintf("Invalid direction '%s'", direction))
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

		params := url.Values{
			"argv1": {strconv.Itoa(ciId1)},
			"argv2": {strconv.Itoa(ciId2)},
			"argv3": {strconv.Itoa(ciRelationTypeId)},
			"argv4": {strconv.Itoa(directionId)},
		}

		jsonRet := CreateCiRelation{}
		err = c.v1.CallWebservice(http.MethodPost, "query", "int_createCiRelation", params, &jsonRet)
		if err != nil {
			err = utilError.FunctionError(err.Error())
			log.Error("Error: ", err)
			return
		}
	}

	return
}

/*

// CreateHistory
// int_createHistory   creates an History-ID
type CreateHistory struct {
	Status string `json:"status"`
}


func (i *Cmdb) CreateHistory() (r CreateHistory, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_createHistory", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// DeleteCi
// int_deleteCi    delete a CI with all dependencies
type DeleteCi struct {
	Status string `json:"status"`
}


func (i *Cmdb) DeleteCi() (r DeleteCi, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_deleteCi", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// DeleteCiAttribute
// int_deleteCiAttribute   delete a ci_attribute-row by id
type DeleteCiAttribute struct {
	Status string `json:"status"`
}


func (i *Cmdb) DeleteCiAttribute() (r DeleteCiAttribute, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_deleteCiAttribute", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}
*/

// DeleteCiRelation
// int_deleteCiRelation    delete a specific ci-relation
type DeleteCiRelation struct {
	Status string `json:"status"`
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

	params := url.Values{
		"argv1": {strconv.Itoa(ciId1)},
		"argv2": {strconv.Itoa(ciId2)},
		"argv3": {strconv.Itoa(ciRelationTypeId)},
	}

	jsonRet := DeleteCiRelation{}
	err = c.v1.CallWebservice(http.MethodPost, "query", "int_deleteCiRelation", params, &jsonRet)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	return
}

/*
// DeleteCiRelationsByCiRelationTypeDirectedFrom
// int_deleteCiRelationsByCiRelationTypeDirectedFrom  deletes all ci-relations with a specific relation-type of a specific CI (direction: from CI)
type DeleteCiRelationsByCiRelationTypeDirectedFrom struct {
	Status string `json:"status"`
}


func (i *Cmdb) DeleteCiRelationsByCiRelationTypeDirectedFrom() (r DeleteCiRelationsByCiRelationTypeDirectedFrom, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_deleteCiRelationsByCiRelationTypeDirectedFrom", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// DeleteCiRelationsByCiRelationTypeDirectedTo
// int_deleteCiRelationsByCiRelationTypeDirectedTo    deletes all ci-relations with a specific relation-type of a specific CI (direction: to CI)
type DeleteCiRelationsByCiRelationTypeDirectedTo struct {
	Status string `json:"status"`
}


func (i *Cmdb) DeleteCiRelationsByCiRelationTypeDirectedTo() (r DeleteCiRelationsByCiRelationTypeDirectedTo, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_deleteCiRelationsByCiRelationTypeDirectedTo", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// DeleteCiRelationsByCiRelationTypeDirectionList
// int_deleteCiRelationsByCiRelationTypeDirectionList     deletes all ci-relations with a specific relation-type of a specific CI
type DeleteCiRelationsByCiRelationTypeDirectionList struct {
	Status string `json:"status"`
}


func (i *Cmdb) DeleteCiRelationsByCiRelationTypeDirectionList() (r DeleteCiRelationsByCiRelationTypeDirectionList, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_deleteCiRelationsByCiRelationTypeDirectionList", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}
*/

// GetAttributeDefaultOption
// int_getAttributeDefaultOption   returns the value of an option
type GetAttributeDefaultOption struct {
	Status string `json:"status"`
	Data   []struct {
		Value string `json:"v"`
	} `json:"data"`
}

func (c *Client) GetAttributeDefaultOption(optionId int) (r string, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetAttributeDefaultOption_" + strconv.Itoa(optionId)
	cached, found := c.v1.Cache.Get(cacheKey)
	if found {
		return cached.(string), nil
	}

	params := url.Values{
		"argv1": {strconv.Itoa(optionId)},
	}

	jsonRet := GetAttributeDefaultOption{}
	err = c.v1.CallWebservice(http.MethodPost, "query", "int_getAttributeDefaultOption", params, &jsonRet)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(strconv.Itoa(optionId) + " - " + v1.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Value
		c.v1.Cache.Set(cacheKey, r, cache.DefaultExpiration)
	default:
		err = utilError.FunctionError(strconv.Itoa(optionId) + " - " + v1.ErrTooManyResults.Error())
	}

	return
}

/*
// GetAttributeDefaultOptionId
// int_getAttributeDefaultOptionId     return the id of a specific attribute and value
type GetAttributeDefaultOptionId struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetAttributeDefaultOptionId() (r GetAttributeDefaultOptionId, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getAttributeDefaultOptionId", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// GetAttributeGroupIdByAttributeGroupName
// int_getAttributeGroupIdByAttributeGroupName     returns the id of an attribute group
type GetAttributeGroupIdByAttributeGroupName struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetAttributeGroupIdByAttributeGroupName() (r GetAttributeGroupIdByAttributeGroupName, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getAttributeGroupIdByAttributeGroupName", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}
*/

// GetAttributeIdByAttributeName
// int_getAttributeIdByAttributeName   returns the id of an attribute
type GetAttributeIdByAttributeNameRet struct {
	Status string       `json:"status"`
	Data   []ResponseId `json:"data"`
}

func (c *Client) GetAttributeIdByAttributeName(name string) (r int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetAttributeIdByAttributeName_" + name
	cached, found := c.v1.Cache.Get(cacheKey)
	if found {
		return cached.(int), nil
	}

	params := url.Values{
		"argv1": {name},
	}

	response := GetAttributeIdByAttributeNameRet{}
	err = c.v1.CallWebservice(http.MethodPost, "query", "int_getAttributeIdByAttributeName", params, &response)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(name + " - " + v1.ErrNoResult.Error())
	case 1:
		r = response.Data[0].Id
		c.v1.Cache.Set(cacheKey, r, cache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v1.ErrTooManyResults.Error())
	}

	return
}

/*
// GetCiAttributeId
// int_getCiAttributeId    returns the id of the first ci_attribute-row with the specific ci_id and attribute_id
type GetCiAttributeId struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetCiAttributeId(ciID int, attrID int) (r GetCiAttributeId, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getCiAttributeId", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}
*/

// GetCiAttributeValue
// int_getCiAttributeValue     get the value of a ci_attribute entry by ci_id and attribute_id
type GetCiAttributeValue struct {
	Status string `json:"status"`
	Data   []struct {
		ID    string `json:"id"`
		Value string `json:"v"`
	} `json:"data"`
}

func (c *Client) GetCiAttributeValue(ciId int, attributeName string, valueType v1.AttributeValueType) (r GetCiAttributeValue, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	attributeId, err := c.GetAttributeIdByAttributeName(attributeName)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	params := url.Values{
		"argv1": {strconv.Itoa(ciId)},
		"argv2": {strconv.Itoa(attributeId)},
		"argv3": {string(valueType)},
	}

	err = c.v1.CallWebservice(http.MethodPost, "query", "int_getCiAttributeValue", params, &r)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	if len(r.Data) == 0 {
		err = utilError.FunctionError(strconv.Itoa(ciId) + ", " + attributeName + " - " + v1.ErrNoResult.Error())
		return
	}

	return
}

func (c *Client) GetCiAttributeValueText(ciId int, attributeName string) (value string, id int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	result, err := c.GetCiAttributeValue(ciId, attributeName, v1.ATTRIBUTE_VALUE_TYPE_TEXT)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	value = result.Data[0].Value
	id, err = strconv.Atoi(result.Data[0].ID)

	return
}

func (c *Client) GetCiAttributeValueDate(ciId int, attributeName string) (value string, id int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	result, err := c.GetCiAttributeValue(ciId, attributeName, v1.ATTRIBUTE_VALUE_TYPE_DATE)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	value = result.Data[0].Value
	id, err = strconv.Atoi(result.Data[0].ID)

	return
}

func (c *Client) GetCiAttributeValueDefault(ciId int, attributeName string) (value string, id int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	result, err := c.GetCiAttributeValue(ciId, attributeName, v1.ATTRIBUTE_VALUE_TYPE_DEFAULT)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	id, err = strconv.Atoi(result.Data[0].ID)

	valueInt, err := strconv.Atoi(result.Data[0].Value)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	value, err = c.GetAttributeDefaultOption(valueInt)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	return
}

func (c *Client) GetCiAttributeValueCi(ciId int, attributeName string) (value string, id int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	result, err := c.GetCiAttributeValue(ciId, attributeName, v1.ATTRIBUTE_VALUE_TYPE_CI)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	value = result.Data[0].Value
	id, err = strconv.Atoi(result.Data[0].ID)

	return
}

/*
// GetCiIdByCiAttributeId
// int_getCiIdByCiAttributeId  returns the ciid of a specific ci_attribute-row
type GetCiIdByCiAttributeId struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetCiIdByCiAttributeId() (r GetCiIdByCiAttributeId, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getCiIdByCiAttributeId", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// GetCiIdByCiAttributeValue
// int_getCiIdByCiAttributeValue   returns the ci_id by a specific attribute_id and value
type GetCiIdByCiAttributeValue struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetCiIdByCiAttributeValue() (r GetCiIdByCiAttributeValue, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getCiIdByCiAttributeValue", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// GetCiProjectMappings
// int_getCiProjectMappings    Get all Projects for a given CI
type GetCiProjectMappings struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetCiProjectMappings() (r GetCiProjectMappings, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getCiProjectMappings", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}
*/

// GetCiRelationCount
// int_getCiRelationCount  returns the number of relations with the given parameters
type GetCiRelationCount struct {
	Status string `json:"status"`
	Data   []struct {
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

	params := url.Values{
		"argv1": {strconv.Itoa(ciId1)},
		"argv2": {strconv.Itoa(ciId2)},
		"argv3": {strconv.Itoa(ciRelationTypeId)},
	}

	jsonRet := GetCiRelationCount{}
	err = c.v1.CallWebservice(http.MethodPost, "query", "int_getCiRelationCount", params, &jsonRet)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return r, err
	}

	errPrefix := strconv.Itoa(ciId1) + ", " + strconv.Itoa(ciId2) + ", " + ciRelationTypeName + "(" + strconv.Itoa(ciRelationTypeId) + ")" + " - "
	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(errPrefix + v1.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Count
	default:
		err = utilError.FunctionError(errPrefix + v1.ErrTooManyResults.Error())
	}

	return

}

// GetCiRelationTypeIdByRelationTypeName
// int_getCiRelationTypeIdByRelationTypeName   returns the id of a relation-type
type GetCiRelationTypeIdByRelationTypeName struct {
	Status string       `json:"status"`
	Data   []ResponseId `json:"data"`
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

	jsonRet := GetCiRelationTypeIdByRelationTypeName{}
	err = c.v2.Query("int_getCiRelationTypeIdByRelationTypeName", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(name + " - " + v1.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Id
		c.v1.Cache.Set(cacheKey, r, cache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v1.ErrTooManyResults.Error())
	}

	return
}

// GetCiTypeIdByCiTypeName
// int_getCiTypeIdByCiTypeName     returns the id for the CI-Type
type GetCiTypeIdByCiTypeName struct {
	Status string       `json:"status"`
	Data   []ResponseId `json:"data"`
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

	params := url.Values{
		"argv1": {name},
	}

	response := GetCiTypeIdByCiTypeName{}
	err = c.v1.CallWebservice(http.MethodPost, "query", "int_getCiTypeIdByCiTypeName", params, &response)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(name + " - " + v1.ErrNoResult.Error())
	case 1:
		r = response.Data[0].Id
		c.v1.Cache.Set(cacheKey, r, cache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v1.ErrTooManyResults.Error())
	}

	return
}

/*
// GetCiTypeOfCi
// int_getCiTypeOfCi   returns the ci-type of a CI
type GetCiTypeOfCi struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetCiTypeOfCi() (r GetCiTypeOfCi, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getCiTypeOfCi", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}
*/

// GetListOfCiIdsByCiRelationDirectionList
// int_getListOfCiIdsByCiRelationDirectionList    returns all related CI-IDs of a specific relation-type
type GetListOfCiIdsByCiRelation struct {
	Status string `json:"status"`
	Data   []struct {
		CiId int `json:"ci_id,string"`
	} `json:"data"`
}

func (c *Client) GetListOfCiIdsByCiRelation(ciId int, ciRelationTypeName string, direction v2.CiRelationDirection) (r []int, err error) {

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

	jsonRet := GetListOfCiIdsByCiRelation{}
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

/*
// GetNumberOfCiAttributes
// int_getNumberOfCiAttributes     returns the number of values for a specific attribute of a CI
type GetNumberOfCiAttributes struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetNumberOfCiAttributes() (r GetNumberOfCiAttributes, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getNumberOfCiAttributes", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// GetProjectIdByProjectName
// int_getProjectIdByProjectName   returns the id of the project with the given name
type GetProjectIdByProjectName struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetProjectIdByProjectName() (r GetProjectIdByProjectName, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getProjectIdByProjectName", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// GetProjects
// int_getProjects     Retrieve all CMDB Projects
type GetProjects struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetProjects() (r GetProjects, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getProjects", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// GetRoleIdByRoleName
// int_getRoleIdByRoleName     returns the id of a role
type GetRoleIdByRoleName struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetRoleIdByRoleName() (r GetRoleIdByRoleName, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getRoleIdByRoleName", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// GetUserIdByUsername
// int_getUserIdByUsername     returns the ID of a infoCMDB-User
type GetUserIdByUsername struct {
	Status string `json:"status"`
}


func (i *Cmdb) GetUserIdByUsername() (r GetUserIdByUsername, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_getUserIdByUsername", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// RemoveCiProjectMapping
// int_removeCiProjectMapping  removes a ci project mapping
type RemoveCiProjectMapping struct {
	Status string `json:"status"`
}


func (i *Cmdb) RemoveCiProjectMapping() (r RemoveCiProjectMapping, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_removeCiProjectMapping", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// SetAttributeRole
// int_setAttributeRole    set permisson for an attribute
type SetAttributeRole struct {
	Status string `json:"status"`
}


func (i *Cmdb) SetAttributeRole() (r SetAttributeRole, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_setAttributeRole", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}

// SetCiTypeOfCi
// int_setCiTypeOfCi   set the ci_type of a CI
type SetCiTypeOfCi struct {
	Status string `json:"status"`
}


func (i *Cmdb) SetCiTypeOfCi() (r SetCiTypeOfCi, err error) {

	if err = i.v2.Login(); err != nil {
		return
	}

	return r, v1.ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.v1.CallWebservice(http.MethodPost,"query", "int_setCiTypeOfCi", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	return
}
*/

type UpdateCiAttribute struct {
	Mode          v2.UpdateMode `json:"mode"`
	Name          string        `json:"name"`
	Value         string        `json:"value"`
	CiAttributeID int           `json:"ciAttributeId"`
	UploadID      string        `json:"uploadId"`
}

type UpdateCiAttributes struct {
	Attributes []UpdateCiAttribute `json:"attributes"`
}

//
type UpdateCiAttributesRequest struct {
	Ci UpdateCiAttributes `json:"ci"`
}

func (c *Client) UpdateCiAttribute(ci int, ua []UpdateCiAttribute) (err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	var errResp clientV2.ResponseError
	resp, err := c.v2.Client.NewRequest().
		SetBody(UpdateCiAttributesRequest{Ci: UpdateCiAttributes{Attributes: ua}}).
		SetAuthToken(c.v2.Config.ApiKey).
		SetError(&errResp).
		Put(fmt.Sprintf("/apiV2/ci/%d", ci))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(errResp.Message + "\n" + errResp.Data)
	}

	return
}

func (c *Client) ListCiByAttributeValue(ci int, ua []UpdateCiAttribute) (err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	var errResp clientV2.ResponseError
	resp, err := c.v2.Client.NewRequest().
		SetError(&errResp).
		SetBody(UpdateCiAttributesRequest{Ci: UpdateCiAttributes{Attributes: ua}}).
		Put(fmt.Sprintf("/apiV2/ci/%d", ci))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(errResp.Message + "\n" + errResp.Data)
	}

	return
}
