package core

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

type ListOfCiIdsOfCiType struct {
	Status string `json:"status"`
	Data   []struct {
		CiID json.Number `json:"ciid"`
	} `json:"data"`
}

func (i *InfoCMDB) GetListOfCiIdsOfCiType(ciTypeID int) (r ListOfCiIdsOfCiType, err error) {
	params := url.Values{
		"argv1": {strconv.Itoa(ciTypeID)},
	}

	ret, err := i.Post("query", "int_getListOfCiIdsOfCiType", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		log.Error(ret)
		return r, err
	}

	return
}

//// Templates for others
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
func (i *InfoCMDB) AddCiProjectMapping(ciID int, projectID int, historyID int) (r AddCiProjectMapping, err error) {
	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
		"argv2": {strconv.Itoa(projectID)},
		"argv3": {strconv.Itoa(historyID)},
	}

	ret, err := i.Post("query", "int_addCiProjectMapping", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		log.Error(ret)
		return r, err
	}

	return
}

//// CreateAttribute
//// int_createAttribute     create an attribute
type CreateAttribute struct {
	Status string `json:"status"`
	CiID   int    `json:"id"`
}

func (i *InfoCMDB) CreateAttribute(ciID int, attrID int) (r CreateAttribute, err error) {
	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
		"argv2": {strconv.Itoa(attrID)},
	}

	ret, err := i.Post("query", "int_createCiAttribute", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		log.Error(ret)
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

func (i *InfoCMDB) CreateCi(ciTypeID int, icon string, historyID int) (r CreateCi, err error) {
	params := url.Values{
		"argv1": {strconv.Itoa(ciTypeID)},
		"argv2": {icon},
		"argv3": {strconv.Itoa(historyID)},
	}

	ret, err := i.Post("query", "int_createCi", params)
	if err != nil {
		log.Error("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Error("Error: ", err)
		log.Error(ret)
		return r, err
	}

	return

}

// GetCi
// int_getCi   Retrieve all informations about a ci
type GetCi struct {
	Status string `json:"status"`
	Data   []struct {
		CiID      json.Number `json:"ci_id"`
		CiTypeID  json.Number `json:"ci_type_id"`
		CiType    string      `json:"ci_type"`
		Project   string      `json:"project"`
		ProjectID json.Number `json:"project_id"`
	} `json:"data"`
}

func (i *InfoCMDB) GetCi(ciID int) (r GetCi, err error) {
	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
	}

	ret, err := i.Post("query", "int_getCi", params)
	if err != nil {
		log.Debugf("Error: %v", err.Error())
		return
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Errorf("Error: %v", err)
		return
	}

	return
}

// GetCiAttributes
// int_getCiAttributes     get all attributes for given ci (:argv1:)

type GetCiAttribute struct {
	CiID                 json.Number `json:"ci_id"`
	CiAttributeID        json.Number `json:"ci_attribute_id"`
	AttributeID          json.Number `json:"attribute_id"`
	AttributeName        string      `json:"attribute_name"`
	AttributeDescription string      `json:"attribute_description"`
	AttributeType        string      `json:"attribute_type"`
	Value                string      `json:"value"`
	ModifiedAt           string      `json:"modified_at"`
}

type GetCiAttributes struct {
	Status string           `json:"status"`
	Data   []GetCiAttribute `json:"data"`
}

func (i *InfoCMDB) GetCiAttributes(ciID int) (r GetCiAttributes, err error) {
	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
	}

	ret, err := i.Post("query", "int_getCiAttributes", params)
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

// CreateCiAttribute
// int_createCiAttribute   creates a ci_attribute-row
type CreateCiAttribute struct {
	Status      string `json:"status"`
	CiID        int    `json:"ci_id"`
	AttributeID int    `json:"attribute_id"`
	HistoryID   int    `json:"history_id"`
}

func (i *InfoCMDB) CreateCiAttribute() (r CreateCiAttribute, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_createCiAttribute", params)
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

// CreateCiRelation
// int_createCiRelation    inserts a relation: argv1 = ci_id_1 argv2 = ci_id_2 argv3 = ci_relation_type_id argv4 = direction
type CreateCiRelation struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) CreateCiRelation() (r CreateCiRelation, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_createCiRelation", params)
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

// CreateHistory
// int_createHistory   creates an History-ID
type CreateHistory struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) CreateHistory() (r CreateHistory, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_createHistory", params)
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

func (i *InfoCMDB) DeleteCi() (r DeleteCi, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_deleteCi", params)
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

func (i *InfoCMDB) DeleteCiAttribute() (r DeleteCiAttribute, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_deleteCiAttribute", params)
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

// DeleteCiRelation
// int_deleteCiRelation    delete a specific ci-relation
type DeleteCiRelation struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) DeleteCiRelation() (r DeleteCiRelation, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_deleteCiRelation", params)
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

// DeleteCiRelationsByCiRelationTypeDirectedFrom
// int_deleteCiRelationsByCiRelationTypeDirectedFrom  deletes all ci-relations with a specific relation-type of a specific CI (direction: from CI)
type DeleteCiRelationsByCiRelationTypeDirectedFrom struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) DeleteCiRelationsByCiRelationTypeDirectedFrom() (r DeleteCiRelationsByCiRelationTypeDirectedFrom, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_deleteCiRelationsByCiRelationTypeDirectedFrom", params)
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

func (i *InfoCMDB) DeleteCiRelationsByCiRelationTypeDirectedTo() (r DeleteCiRelationsByCiRelationTypeDirectedTo, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_deleteCiRelationsByCiRelationTypeDirectedTo", params)
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

func (i *InfoCMDB) DeleteCiRelationsByCiRelationTypeDirectionList() (r DeleteCiRelationsByCiRelationTypeDirectionList, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_deleteCiRelationsByCiRelationTypeDirectionList", params)
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

// GetAttributeDefaultOption
// int_getAttributeDefaultOption   returns the value of an option
type GetAttributeDefaultOption struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetAttributeDefaultOption() (r GetAttributeDefaultOption, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getAttributeDefaultOption", params)
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

// GetAttributeDefaultOptionId
// int_getAttributeDefaultOptionId     return the id of a specific attribute and value
type GetAttributeDefaultOptionId struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetAttributeDefaultOptionId() (r GetAttributeDefaultOptionId, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getAttributeDefaultOptionId", params)
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

func (i *InfoCMDB) GetAttributeGroupIdByAttributeGroupName() (r GetAttributeGroupIdByAttributeGroupName, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getAttributeGroupIdByAttributeGroupName", params)
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

// GetAttributeIdByAttributeName
// int_getAttributeIdByAttributeName   returns the id of an attribute
type GetAttributeIdByAttributeName struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetAttributeIdByAttributeName() (r GetAttributeIdByAttributeName, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getAttributeIdByAttributeName", params)
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

// GetCiAttributeId
// int_getCiAttributeId    returns the id of the first ci_attribute-row with the specific ci_id and attribute_id
type GetCiAttributeId struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetCiAttributeId(ciID int, attrID int) (r GetCiAttributeId, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiAttributeId", params)
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

// GetCiAttributeValue
// int_getCiAttributeValue     get the value of a ci_attribute entry by ci_id and attribute_id
type GetCiAttributeValue struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetCiAttributeValue() (r GetCiAttributeValue, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiAttributeValue", params)
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

// GetCiIdByCiAttributeId
// int_getCiIdByCiAttributeId  returns the ciid of a specific ci_attribute-row
type GetCiIdByCiAttributeId struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetCiIdByCiAttributeId() (r GetCiIdByCiAttributeId, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiIdByCiAttributeId", params)
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

func (i *InfoCMDB) GetCiIdByCiAttributeValue() (r GetCiIdByCiAttributeValue, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiIdByCiAttributeValue", params)
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

func (i *InfoCMDB) GetCiProjectMappings() (r GetCiProjectMappings, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiProjectMappings", params)
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

// GetCiRelationCount
// int_getCiRelationCount  returns the number of relations with the given parameters
type GetCiRelationCount struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetCiRelationCount() (r GetCiRelationCount, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiRelationCount", params)
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

// GetCiRelationTypeIdByRelationTypeName
// int_getCiRelationTypeIdByRelationTypeName   returns the id of a relation-type
type GetCiRelationTypeIdByRelationTypeName struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetCiRelationTypeIdByRelationTypeName() (r GetCiRelationTypeIdByRelationTypeName, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiRelationTypeIdByRelationTypeName", params)
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

// GetCiTypeIdByCiTypeName
// int_getCiTypeIdByCiTypeName     returns the id for the CI-Type
type GetCiTypeIdByCiTypeName struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetCiTypeIdByCiTypeName() (r GetCiTypeIdByCiTypeName, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiTypeIdByCiTypeName", params)
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

// GetCiTypeOfCi
// int_getCiTypeOfCi   returns the ci-type of a CI
type GetCiTypeOfCi struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetCiTypeOfCi() (r GetCiTypeOfCi, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getCiTypeOfCi", params)
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

// GetListOfCiIdsByCiRelationDirectedFrom
// int_getListOfCiIdsByCiRelationDirectedFrom     returns all related CI-IDs of a specific relation-type (direction: from CI)
type GetListOfCiIdsByCiRelationDirectedFrom struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetListOfCiIdsByCiRelationDirectedFrom() (r GetListOfCiIdsByCiRelationDirectedFrom, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getListOfCiIdsByCiRelationDirectedFrom", params)
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

// GetListOfCiIdsByCiRelationDirectedTo
// int_getListOfCiIdsByCiRelationDirectedTo   returns all related CI-IDs of a specific relation-type (direction: to CI)
type GetListOfCiIdsByCiRelationDirectedTo struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetListOfCiIdsByCiRelationDirectedTo() (r GetListOfCiIdsByCiRelationDirectedTo, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getListOfCiIdsByCiRelationDirectedTo", params)
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

// GetListOfCiIdsByCiRelationDirectionList
// int_getListOfCiIdsByCiRelationDirectionList    returns all related CI-IDs of a specific relation-type
type GetListOfCiIdsByCiRelationDirectionList struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetListOfCiIdsByCiRelationDirectionList() (r GetListOfCiIdsByCiRelationDirectionList, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getListOfCiIdsByCiRelationDirectionList", params)
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

// GetNumberOfCiAttributes
// int_getNumberOfCiAttributes     returns the number of values for a specific attribute of a CI
type GetNumberOfCiAttributes struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) GetNumberOfCiAttributes() (r GetNumberOfCiAttributes, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getNumberOfCiAttributes", params)
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

func (i *InfoCMDB) GetProjectIdByProjectName() (r GetProjectIdByProjectName, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getProjectIdByProjectName", params)
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

func (i *InfoCMDB) GetProjects() (r GetProjects, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getProjects", params)
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

func (i *InfoCMDB) GetRoleIdByRoleName() (r GetRoleIdByRoleName, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getRoleIdByRoleName", params)
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

func (i *InfoCMDB) GetUserIdByUsername() (r GetUserIdByUsername, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_getUserIdByUsername", params)
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

func (i *InfoCMDB) RemoveCiProjectMapping() (r RemoveCiProjectMapping, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_removeCiProjectMapping", params)
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

func (i *InfoCMDB) SetAttributeRole() (r SetAttributeRole, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_setAttributeRole", params)
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

func (i *InfoCMDB) SetCiTypeOfCi() (r SetCiTypeOfCi, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_setCiTypeOfCi", params)
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

// UpdateCiAttribute
// int_updateCiAttribute   updates a specific ci_attribute_row argv1 = ci_attribute-ID argv2 = column argv3 = value argv4 = history_id
type UpdateCiAttribute struct {
	Status string `json:"status"`
}

func (i *InfoCMDB) UpdateCiAttribute() (r UpdateCiAttribute, err error) {
	return r, ErrNotImplemented // TODO FIXME
	params := url.Values{
		// "argv1": {strconv.Itoa(%PARAM1%)},
		// "argv2": {strconv.Itoa(%PARAM2%)},
		// "argv3": {strconv.Itoa(%PARAM3%)},
		// "argv4": {strconv.Itoa(%PARAM4%)},
	}

	ret, err := i.Post("query", "int_updateCiAttribute", params)
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
