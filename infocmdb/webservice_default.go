package infocmdb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// QueryWebservices allows you to call a generic webservice(arg1: ws) with the providing params
// Return: json string
func (c *Client) QueryWebservice(ws string, params map[string]string) (resp string, err error) {
	log.Debugf("Querying webservice %v with params %v", ws, params)

	if err = c.v2.Login(); err != nil {
		return
	}

	resp, err = c.v2.QueryRaw(ws, params)
	if err != nil {
		log.Error("Error: ", err)
	}

	log.Debugf("Result: %v", resp)
	return
}

// Query allows you to call a generic webservice(arg1: ws) with the providing params and a reference
// to a result. It will take the built in resty function to deserialize the result
// Return: error
func (c *Client) Query(ws string, out interface{}, params map[string]string) (err error) {
	log.Debugf("Querying webservice %v with params %v", ws, params)

	if err = c.v2.Login(); err != nil {
		return
	}

	if err = c.v2.Query(ws, out, params); err != nil {
		log.Error("Error: ", err)
	}

	log.Debugf("Result: %v", out)
	return
}

type CiIds []int

type getListOfCiIdsOfCiType struct {
	Status string `json:"status"`
	Data   []struct {
		CiID int `json:"ciid,string"`
	} `json:"data"`
}

type responseId struct {
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

type getListOfCiIdsByAttributeValue struct {
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

// AddCiProjectMapping
// int_addCiProjectMapping     add project-mapping to a ci
//
// insert into ci_project (ci_id, project_id, history_id)
// select :argv1:, :argv2:, :argv3:
// from dual
// where not exists(select id from ci_project where ci_id = :argv1: and project_id = :argv2:)

func (c *Client) AddCiProjectMapping(ciID int, projectID int, historyID int) (err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciID),
		"argv2": strconv.Itoa(projectID),
		"argv3": strconv.Itoa(historyID),
	}

	err = c.v2.Query("int_addCiProjectMapping", nil, params)
	if err != nil {
		log.Error("Error: ", err)
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

	params := map[string]string{
		"argv1": strconv.Itoa(ciID),
		"argv2": strconv.Itoa(attrID),
	}

	err = c.v2.Query("int_createCiAttribute", &r, params)
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

	params := map[string]string{
		"argv1": strconv.Itoa(ciTypeID),
		"argv2": icon,
		"argv3": strconv.Itoa(historyID),
	}

	err = c.v2.Query("int_createCi", &r, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return r, err
	}
	return
}

// getCi
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
type getCi struct {
	Status string `json:"status"`
	Data   []struct {
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

// getCiAttributes
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

type getCiAttributes struct {
	Status string        `json:"status"`
	Data   []CiAttribute `json:"data"`
}

func (c *Client) GetCiAttributes(ciID int) (r []CiAttribute, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciID),
	}

	jsonRet := getCiAttributes{}
	err = c.v2.Query("int_getCiAttributes", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
	}
	r = jsonRet.Data
	return
}

// createCiRelation
// int_createCiRelation    inserts a relation: argv1 = ci_id_1 argv2 = ci_id_2 argv3 = ci_relation_type_id argv4 = direction
type createCiRelation struct {
	Status string `json:"status"`
}

func (c *Client) CreateCiRelation(ciId1 int, ciId2 int, ciRelationTypeName string, direction v2.CiRelationDirection) (err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	var directionId int
	switch direction {
	case v2.CI_RELATION_DIRECTION_DIRECTED_FROM:
		directionId = 1
	case v2.CI_RELATION_DIRECTION_DIRECTED_TO:
		directionId = 2
	case v2.CI_RELATION_DIRECTION_BIDIRECTIONAL:
		directionId = 3
	case v2.CI_RELATION_DIRECTION_OMNIDIRECTIONAL:
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

// deleteCiRelation
// int_deleteCiRelation    delete a specific ci-relation
type deleteCiRelation struct {
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

// getAttributeDefaultOption
// int_getAttributeDefaultOption   returns the value of an option
type getAttributeDefaultOption struct {
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

	params := map[string]string{
		"argv1": strconv.Itoa(optionId),
	}

	jsonRet := getAttributeDefaultOption{}
	err = c.v2.Query("int_getAttributeDefaultOption", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(jsonRet.Data) {
	case 0:
		err = utilError.FunctionError(strconv.Itoa(optionId) + " - " + v2.ErrNoResult.Error())
	case 1:
		r = jsonRet.Data[0].Value
		c.v1.Cache.Set(cacheKey, r, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(strconv.Itoa(optionId) + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

// GetAttributeIdByAttributeName
// int_getAttributeIdByAttributeName   returns the id of an attribute
type getAttributeIdByAttributeNameRet struct {
	Status string       `json:"status"`
	Data   []responseId `json:"data"`
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

	params := map[string]string{
		"argv1": name,
	}

	response := getAttributeIdByAttributeNameRet{}
	err = c.v2.Query("int_getAttributeIdByAttributeName", &response, params)
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

// GetCiAttributeValue
// int_getCiAttributeValue     get the value of a ci_attribute entry by ci_id and attribute_id
type GetCiAttributeValue struct {
	Status string `json:"status"`
	Data   []struct {
		ID    string `json:"id"`
		Value string `json:"v"`
	} `json:"data"`
}

func (c *Client) GetCiAttributeValue(ciId int, attributeName string, valueType v2.AttributeValueType) (r GetCiAttributeValue, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	attributeId, err := c.GetAttributeIdByAttributeName(attributeName)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(ciId),
		"argv2": strconv.Itoa(attributeId),
		"argv3": string(valueType),
	}

	err = c.v2.Query("int_getCiAttributeValue", &r, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	if len(r.Data) == 0 {
		err = utilError.FunctionError(strconv.Itoa(ciId) + ", " + attributeName + " - " + v2.ErrNoResult.Error())
		return
	}

	return
}

func (c *Client) GetCiAttributeValueText(ciId int, attributeName string) (value string, id int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	result, err := c.GetCiAttributeValue(ciId, attributeName, v2.ATTRIBUTE_VALUE_TYPE_TEXT)
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

	result, err := c.GetCiAttributeValue(ciId, attributeName, v2.ATTRIBUTE_VALUE_TYPE_DATE)
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

	result, err := c.GetCiAttributeValue(ciId, attributeName, v2.ATTRIBUTE_VALUE_TYPE_DEFAULT)
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

	result, err := c.GetCiAttributeValue(ciId, attributeName, v2.ATTRIBUTE_VALUE_TYPE_CI)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		return
	}

	value = result.Data[0].Value
	id, err = strconv.Atoi(result.Data[0].ID)

	return
}

// getCiRelationCount
// int_getCiRelationCount  returns the number of relations with the given parameters
type getCiRelationCount struct {
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

// getCiRelationTypeIdByRelationTypeName
// int_getCiRelationTypeIdByRelationTypeName   returns the id of a relation-type
type getCiRelationTypeIdByRelationTypeName struct {
	Status string       `json:"status"`
	Data   []responseId `json:"data"`
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

// getCiTypeIdByCiTypeName
// int_getCiTypeIdByCiTypeName     returns the id for the CI-Type
type getCiTypeIdByCiTypeName struct {
	Status string       `json:"status"`
	Data   []responseId `json:"data"`
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

// GetListOfCiIdsByCiRelationDirectionList
// int_getListOfCiIdsByCiRelationDirectionList    returns all related CI-IDs of a specific relation-type
type getListOfCiIdsByCiRelation struct {
	Status string `json:"status"`
	Data   []struct {
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

func (c *Client) UpdateCiAttribute(ci int, ua []v2.UpdateCiAttribute) (err error) {
	return c.v2.UpdateCiAttribute(ci, ua)
}

func (c *Client) GetWorkflowContext(workflowInstanceId int) (workflowContext *v2.WorkflowContext, err error) {
	return c.v2.GetWorkflowContext(workflowInstanceId)
}
