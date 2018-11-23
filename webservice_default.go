package infocmdblibrary

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

type ListOfCiIdsOfCiType struct {
	Status string `json:"status"`
	Data   []struct {
		CiID json.Number `json:"ciid"`
	} `json:"data"`
}

func (i *InfoCmdbGoLib) GetListOfCiIdsOfCiType(ciTypeID int) (ListOfCiIdsOfCiType, error) {
	r := ListOfCiIdsOfCiType{}

	params := url.Values{
		"argv1": {strconv.Itoa(ciTypeID)},
	}

	ret, err := i.WS.client.Post("query", "int_getListOfCiIdsOfCiType", params)
	if err != nil {
		log.Println("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Println("Error: ", err)
		log.Println(ret)
		return r, err
	}

	return r, nil
}

//// Templates for others
//type AddCiProjectMapping struct {
//}
//
//// AddCiProjectMapping
//// int_addCiProjectMapping     add project-mapping to a ci         Aktiv
////
//// insert into ci_project (ci_id, project_id, history_id)
//// select :argv1:, :argv2:, :argv3:
//// from dual
//// where not exists(select id from ci_project where ci_id = :argv1: and project_id = :argv2:)
////func (i *InfoCmdbGoLib) AddCiProjectMapping(ciID int, projectID int, historyID int) (AddCiProjectMapping, error) {
////
////}
//// CreateAttribute
//// int_createAttribute     create an attribute         Aktiv
//type CreateAttribute struct {
//}
//
//func (i *InfoCmdbGoLib) CreateAttribute() (CreateAttribute, error) {
//
//}
//// CreateAttributeGroup
//// int_createAttributeGroup    create an attribute-group       Aktiv
//type CreateAttributeGroup struct {
//}
//
//func (i *InfoCmdbGoLib) CreateAttributeGroup() (CreateAttributeGroup, error) {
//
//}
//// CreateCi
//// int_createCi    create a CI         Aktiv
//type CreateCi struct {
//}
//
//func (i *InfoCmdbGoLib) CreateCi() (CreateCi, error) {
//
//}
//// CreateCiAttribute
//// int_createCiAttribute   creates a ci_attribute-row      Aktiv
//type CreateCiAttribute struct {
//}
//
//func (i *InfoCmdbGoLib) CreateCiAttribute() (CreateCiAttribute, error) {
//
//}
//// CreateCiRelation
//// int_createCiRelation    inserts a relation: argv1 = ci_id_1 argv2 = ci_id_2 argv3 = ci_relation_type_id argv4 = direction       Aktiv
//type CreateCiRelation struct {
//}
//
//func (i *InfoCmdbGoLib) CreateCiRelation() (CreateCiRelation, error) {
//
//}
//// CreateHistory
//// int_createHistory   creates an History-ID       Aktiv
//type CreateHistory struct {
//}
//
//func (i *InfoCmdbGoLib) CreateHistory() (CreateHistory, error) {
//
//}
//// DeleteCi
//// int_deleteCi    delete a CI with all dependencies       Aktiv
//type DeleteCi struct {
//}
//
//func (i *InfoCmdbGoLib) DeleteCi() (DeleteCi, error) {
//
//}
//// DeleteCiAttribute
//// int_deleteCiAttribute   delete a ci_attribute-row by id         Aktiv
//type DeleteCiAttribute struct {
//}
//
//func (i *InfoCmdbGoLib) DeleteCiAttribute() (DeleteCiAttribute, error) {
//
//}
//// DeleteCiRelation
//// int_deleteCiRelation    delete a specific ci-relation       Aktiv
//type DeleteCiRelation struct {
//}
//
//func (i *InfoCmdbGoLib) DeleteCiRelation() (DeleteCiRelation, error) {
//
//}
//// DeleteCiRelationsByCiRelationType_directedFrom
//// int_deleteCiRelationsByCiRelationType_directedFrom  deletes all ci-relations with a specific relation-type of a specific CI (direction: from CI)        Aktiv
//type DeleteCiRelationsByCiRelationType_directedFrom struct {
//}
//
//func (i *InfoCmdbGoLib) DeleteCiRelationsByCiRelationType_directedFrom() (DeleteCiRelationsByCiRelationType_directedFrom, error) {
//
//}
//// DeleteCiRelationsByCiRelationType_directedTo
//// int_deleteCiRelationsByCiRelationType_directedTo    deletes all ci-relations with a specific relation-type of a specific CI (direction: to CI)      Aktiv
//type DeleteCiRelationsByCiRelationType_directedTo struct {
//}
//
//func (i *InfoCmdbGoLib) DeleteCiRelationsByCiRelationType_directedTo() (DeleteCiRelationsByCiRelationType_directedTo, error) {
//
//}
//// DeleteCiRelationsByCiRelationType_directionList
//// int_deleteCiRelationsByCiRelationType_directionList     deletes all ci-relations with a specific relation-type of a specific CI         Aktiv
//type DeleteCiRelationsByCiRelationType_directionList struct {
//}
//
//func (i *InfoCmdbGoLib) DeleteCiRelationsByCiRelationType_directionList() (DeleteCiRelationsByCiRelationType_directionList, error) {
//
//}
//// GetAttributeDefaultOption
//// int_getAttributeDefaultOption   returns the value of an option      Aktiv
//type GetAttributeDefaultOption struct {
//}
//
//func (i *InfoCmdbGoLib) GetAttributeDefaultOption() (GetAttributeDefaultOption, error) {
//
//}
//// GetAttributeDefaultOptionId
//// int_getAttributeDefaultOptionId     return the id of a specific attribute and value         Aktiv
//type GetAttributeDefaultOptionId struct {
//}
//
//func (i *InfoCmdbGoLib) GetAttributeDefaultOptionId() (GetAttributeDefaultOptionId, error) {
//
//}
//// GetAttributeGroupIdByAttributeGroupName
//// int_getAttributeGroupIdByAttributeGroupName     returns the id of an attribute group        Aktiv
//type GetAttributeGroupIdByAttributeGroupName struct {
//}
//
//func (i *InfoCmdbGoLib) GetAttributeGroupIdByAttributeGroupName() (GetAttributeGroupIdByAttributeGroupName, error) {
//
//}
//// GetAttributeIdByAttributeName
//// int_getAttributeIdByAttributeName   returns the id of an attribute      Aktiv
//type GetAttributeIdByAttributeName struct {
//}
//
//func (i *InfoCmdbGoLib) GetAttributeIdByAttributeName() (GetAttributeIdByAttributeName, error) {
//
//}
// GetCi
// int_getCi   Retrieve all informations about a ci        Aktiv
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

func (i *InfoCmdbGoLib) GetCi(ciID int) (GetCi, error) {

	r := GetCi{}

	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
	}

	ret, err := i.WS.client.Post("query", "int_getCi", params)
	if err != nil {
		log.Println("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Println("Error: ", err)
		return r, err
	}

	return r, nil
}

//// GetCiAttributeId
//// int_getCiAttributeId    returns the id of the first ci_attribute-row with the specific ci_id and attribute_id       Aktiv
//type GetCiAttributeId struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiAttributeId() (GetCiAttributeId, error) {
//
//}
// GetCiAttributes
// int_getCiAttributes     get all attributes for given ci (:argv1:)       Aktiv

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

func (i *InfoCmdbGoLib) GetCiAttributes(ciID int) (GetCiAttributes, error) {

	r := GetCiAttributes{}

	params := url.Values{
		"argv1": {strconv.Itoa(ciID)},
	}

	ret, err := i.WS.client.Post("query", "int_getCiAttributes", params)
	if err != nil {
		log.Println("Error: ", err)
		return r, err
	}

	err = json.Unmarshal([]byte(ret), &r)
	if err != nil {
		log.Println("Error: ", err)
		return r, err
	}

	return r, nil
}

//// GetCiAttributeValue
//// int_getCiAttributeValue     get the value of a ci_attribute entry by ci_id and attribute_id         Aktiv
//type GetCiAttributeValue struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiAttributeValue() (GetCiAttributeValue, error) {
//
//}
//// GetCiIdByCiAttributeId
//// int_getCiIdByCiAttributeId  returns the ciid of a specific ci_attribute-row         Aktiv
//type GetCiIdByCiAttributeId struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiIdByCiAttributeId() (GetCiIdByCiAttributeId, error) {
//
//}
//// GetCiIdByCiAttributeValue
//// int_getCiIdByCiAttributeValue   returns the ci_id by a specific attribute_id and value      Aktiv
//type GetCiIdByCiAttributeValue struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiIdByCiAttributeValue() (GetCiIdByCiAttributeValue, error) {
//
//}
//// GetCiProjectMappings
//// int_getCiProjectMappings    Get all Projects for a given CI         Aktiv
//type GetCiProjectMappings struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiProjectMappings() (GetCiProjectMappings, error) {
//
//}
//// GetCiRelationCount
//// int_getCiRelationCount  returns the number of relations with the given parameters       Aktiv
//type GetCiRelationCount struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiRelationCount() (GetCiRelationCount, error) {
//
//}
//// GetCiRelationTypeIdByRelationTypeName
//// int_getCiRelationTypeIdByRelationTypeName   returns the id of a relation-type       Aktiv
//type GetCiRelationTypeIdByRelationTypeName struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiRelationTypeIdByRelationTypeName() (GetCiRelationTypeIdByRelationTypeName, error) {
//
//}
//// GetCiTypeIdByCiTypeName
//// int_getCiTypeIdByCiTypeName     returns the id for the CI-Type      Aktiv
//type GetCiTypeIdByCiTypeName struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiTypeIdByCiTypeName() (GetCiTypeIdByCiTypeName, error) {
//
//}
//// GetCiTypeOfCi
//// int_getCiTypeOfCi   returns the ci-type of a CI         Aktiv
//type GetCiTypeOfCi struct {
//}
//
//func (i *InfoCmdbGoLib) GetCiTypeOfCi() (GetCiTypeOfCi, error) {
//
//}
//// GetListOfCiIdsByCiRelation_directedFrom
//// int_getListOfCiIdsByCiRelation_directedFrom     returns all related CI-IDs of a specific relation-type (direction: from CI)         Aktiv
//type GetListOfCiIdsByCiRelation_directedFrom struct {
//}
//
//func (i *InfoCmdbGoLib) GetListOfCiIdsByCiRelation_directedFrom() (GetListOfCiIdsByCiRelation_directedFrom, error) {
//
//}
//// GetListOfCiIdsByCiRelation_directedTo
//// int_getListOfCiIdsByCiRelation_directedTo   returns all related CI-IDs of a specific relation-type (direction: to CI)       Aktiv
//type GetListOfCiIdsByCiRelation_directedTo struct {
//}
//
//func (i *InfoCmdbGoLib) GetListOfCiIdsByCiRelation_directedTo() (GetListOfCiIdsByCiRelation_directedTo, error) {
//
//}
//// GetListOfCiIdsByCiRelation_directionList
//// int_getListOfCiIdsByCiRelation_directionList    returns all related CI-IDs of a specific relation-type      Aktiv
//type GetListOfCiIdsByCiRelation_directionList struct {
//}
//
//func (i *InfoCmdbGoLib) GetListOfCiIdsByCiRelation_directionList() (GetListOfCiIdsByCiRelation_directionList, error) {
//
//}
//
//// GetNumberOfCiAttributes
//// int_getNumberOfCiAttributes     returns the number of values for a specific attribute of a CI       Aktiv
//type GetNumberOfCiAttributes struct {
//}
//
//func (i *InfoCmdbGoLib) GetNumberOfCiAttributes() (GetNumberOfCiAttributes, error) {
//
//}
//// GetProjectIdByProjectName
//// int_getProjectIdByProjectName   returns the id of the project with the given name       Aktiv
//type GetProjectIdByProjectName struct {
//}
//
//func (i *InfoCmdbGoLib) GetProjectIdByProjectName() (GetProjectIdByProjectName, error) {
//
//}
//// GetProjects
//// int_getProjects     Retrieve all CMDB Projects      Aktiv
//type GetProjects struct {
//}
//
//func (i *InfoCmdbGoLib) GetProjects() (GetProjects, error) {
//
//}
//// GetRoleIdByRoleName
//// int_getRoleIdByRoleName     returns the id of a role        Aktiv
//type GetRoleIdByRoleName struct {
//}
//
//func (i *InfoCmdbGoLib) GetRoleIdByRoleName() (GetRoleIdByRoleName, error) {
//
//}
//// GetUserIdByUsername
//// int_getUserIdByUsername     returns the ID of a infoCMDB-User       Aktiv
//type GetUserIdByUsername struct {
//}
//
//func (i *InfoCmdbGoLib) GetUserIdByUsername() (GetUserIdByUsername, error) {
//
//}
//// RemoveCiProjectMapping
//// int_removeCiProjectMapping  removes a ci project mapping        Aktiv
//type RemoveCiProjectMapping struct {
//}
//
//func (i *InfoCmdbGoLib) RemoveCiProjectMapping() (RemoveCiProjectMapping, error) {
//
//}
//// SetAttributeRole
//// int_setAttributeRole    set permisson for an attribute      Aktiv
//type SetAttributeRole struct {
//}
//
//func (i *InfoCmdbGoLib) SetAttributeRole() (SetAttributeRole, error) {
//
//}
//// SetCiTypeOfCi
//// int_setCiTypeOfCi   set the ci_type of a CI         Aktiv
//type SetCiTypeOfCi struct {
//}
//
//func (i *InfoCmdbGoLib) SetCiTypeOfCi() (SetCiTypeOfCi, error) {
//
//}
//// UpdateCiAttribute
//// int_updateCiAttribute   updates a specific ci_attribute_row argv1 = ci_attribute-ID argv2 = column argv3 = value argv4 = history_id         Akti
//type UpdateCiAttribute struct {
//}
//
//func (i *InfoCmdbGoLib) UpdateCiAttribute() (UpdateCiAttribute, error) {
//
//}
