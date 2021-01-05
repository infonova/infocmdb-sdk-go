package infocmdb

import (
	"errors"
	"strconv"
	"strings"

	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
)

type CiAttributes = []CiAttribute

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
	Data []CiAttribute `json:"data"`
}

func (c *Client) GetCiAttributes(ciId int) (ciAttributes CiAttributes, err error) {
	ciIdToAttributesMap, err := c.GetMapOfCiAttributes([]int{ciId})
	if err != nil {
		return
	}

	ciAttributes = ciIdToAttributesMap[ciId]
	return
}

func (c *Client) GetMapOfCiAttributes(ciIds []int) (ciIdToAttributesMap map[int]CiAttributes, err error) {
	ciIdToAttributesMap = map[int]CiAttributes{}

	if len(ciIds) == 0 {
		return
	}

	if err = c.v2.Login(); err != nil {
		return
	}

	commaSeparatedCiIds := ""
	for _, ciId := range ciIds {
		if commaSeparatedCiIds != "" {
			commaSeparatedCiIds += ", "
		}

		commaSeparatedCiIds += strconv.Itoa(ciId)
	}

	params := map[string]string{
		"argv1": commaSeparatedCiIds,
	}

	jsonRet := getCiAttributes{}
	err = c.v2.Query("int_getCiAttributes", &jsonRet, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
	}

	for _, ciAttribute := range jsonRet.Data {
		ciAttributes, ok := ciIdToAttributesMap[ciAttribute.CiID]
		if !ok {
			ciAttributes = CiAttributes{}
		}
		ciAttributes = append(ciAttributes, ciAttribute)
		ciIdToAttributesMap[ciAttribute.CiID] = ciAttributes
	}

	return
}

type getAttributeDefaultOption struct {
	Data []struct {
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

type getAttributeDefaultOptionId struct {
	Data []responseId `json:"data"`
}

func (c *Client) GetAttrDefaultOptionIdByAttrId(attrId int, optionValue string) (attrDefaultOptionId int, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	attrIdString := strconv.Itoa(attrId)

	cacheKey := "GetAttrDefaultOptionIdByAttrId_" + attrIdString + "_" + optionValue
	cached, found := c.v2.Cache.Get(cacheKey)
	if found {
		return cached.(int), nil
	}

	params := map[string]string{
		"argv1": attrIdString,
		"argv2": optionValue,
	}

	response := getAttributeDefaultOptionId{}
	err = c.v2.Query("int_getAttributeDefaultOptionId", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(attrIdString + ", " + optionValue + " - " + v2.ErrNoResult.Error())
	case 1:
		attrDefaultOptionId = response.Data[0].Id
		c.v2.Cache.Set(cacheKey, attrDefaultOptionId, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(attrIdString + ", " + optionValue + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

func (c *Client) GetAttrDefaultOptionIdByAttrName(attrName string, optionValue string) (attrDefaultOptionId int, err error) {
	attrId, err := c.GetAttributeIdByAttributeName(attrName)
	if err != nil {
		return
	}

	attrDefaultOptionId, err = c.GetAttrDefaultOptionIdByAttrId(attrId, optionValue)
	return
}

type getAttributeIdByAttributeNameRet struct {
	Data []responseId `json:"data"`
}

func (c *Client) GetAttributeIdByAttributeName(name string) (attrId int, err error) {
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
		attrId = response.Data[0].Id
		c.v1.Cache.Set(cacheKey, attrId, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

type GetCiAttributeValue struct {
	Data []struct {
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

func (c *Client) UpdateCiAttribute(ci int, ua []v2.UpdateCiAttribute) (err error) {
	return c.v2.UpdateCiAttribute(ci, ua)
}

type AttributeType int

const (
	AT_INPUT AttributeType = iota + 1
	AT_TEXTAREA
	AT_TEXTEDIT
	AT_SELECTFIELD
	AT_CHECKBOX
	AT_RADIO
	AT_DATE
	AT_DATETIME
	AT_ZAHLUNGSMITTEL
	AT_PASSWORD
	AT_LINK
	AT_ATTACHMENT
	AT_SCRIPT
	AT_EXECUTEABLE
	AT_QUERY
	AT_CITYPE
	AT_INFO
	AT_QUERYPERSIST
	AT_CITYPEPERSIST
	AT_FILTER
	AT_SELECTQUERY
	AT_SELECTPOPUP
)

type Columns int

const (
	C_ONE Columns = iota + 1
	C_TWO
)

type Multiselect int

const (
	MS_ZERO Multiselect = iota
	MS_ONE
	MS_TWO
)

type respCreateAttribute struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []responseId `json:"data"`
}

type AttributeParams struct {
	Name                string
	Description         string
	Note                string
	Hint                string
	AttributeType       AttributeType
	AttributeGroupName  string
	OrderNumber         int
	Column              Columns //values in db enum('1','2')
	IsUnique            bool
	IsNumeric           bool
	IsBold              bool
	IsEvent             bool
	IsUniqueCheck       bool
	IsAutocomplete      bool
	IsMultiselect       Multiselect //values in db enum('0','1','2')
	IsProjectRestricted bool
	Regex               string
	ScriptName          string
	InputMaxlength      int
	TextareaCols        int
	TextareaRows        int
	IsActive            bool
	Historicize         bool
	UserId              int
}

func (c *Client) NewAttributeParams() (params *AttributeParams) {
	params = &AttributeParams{
		AttributeType: 0,
		Column:        C_ONE,
		IsMultiselect: MS_ZERO,
		IsActive:      false,
		Historicize:   true,
		UserId:        0,
	}
	return
}

func (c *Client) CreateAttribute(attributeParams *AttributeParams) (attributeId int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	existingAttributeId, err := c.GetAttributeIdByAttributeName(attributeParams.Name)
	if err != nil && strings.Contains(err.Error(), "query returned no result") == false {
		return 0, err
	}

	if existingAttributeId == 0 {

		columns := []string{
			"name",
			"description",
			"note",
			"hint",
			"attribute_type_id",
			"attribute_group_id",
			"order_number",
			"column",
			"is_unique",
			"is_numeric",
			"is_bold",
			"is_event",
			"is_unique_check",
			"is_autocomplete",
			"is_multiselect",
			"is_project_restricted",
			"regex",
			"script_name",
			"input_maxlength",
			"textarea_cols",
			"textarea_rows",
			"is_active",
			"user_id",
			"historicize",
		}

		attributeGroupId, err := c.GetAttributeGroupIdByName(attributeParams.AttributeGroupName)
		if err != nil {
			return 0, err
		}

		values := []string{
			attributeParams.Name,
			attributeParams.Description,
			attributeParams.Note,
			attributeParams.Hint,
			strconv.Itoa(int(attributeParams.AttributeType)),
			strconv.Itoa(attributeGroupId),
			strconv.Itoa(attributeParams.OrderNumber),
			strconv.Itoa(int(attributeParams.Column)),
			convertBoolToString[attributeParams.IsUnique],
			convertBoolToString[attributeParams.IsNumeric],
			convertBoolToString[attributeParams.IsBold],
			convertBoolToString[attributeParams.IsEvent],
			convertBoolToString[attributeParams.IsUniqueCheck],
			convertBoolToString[attributeParams.IsAutocomplete],
			strconv.Itoa(int(attributeParams.IsMultiselect)),
			convertBoolToString[attributeParams.IsProjectRestricted],
			attributeParams.Regex,
			attributeParams.ScriptName,
			strconv.Itoa(attributeParams.InputMaxlength),
			strconv.Itoa(attributeParams.TextareaCols),
			strconv.Itoa(attributeParams.TextareaRows),
			convertBoolToString[attributeParams.IsActive],
			strconv.Itoa(attributeParams.UserId),
			convertBoolToString[attributeParams.Historicize],
		}

		params := map[string]string{
			"argv1": "`" + strings.Join(columns, "`, `") + "`",
			"argv2": "'" + strings.Join(values, "', '") + "'",
		}

		response := respCreateAttribute{}
		err = c.v2.Query("int_createAttribute", &response, params)
		if err != nil {
			err = utilError.FunctionError(err.Error())
			log.Error("Error: ", err)
			return 0, err
		}

		switch len(response.Data) {
		case 0:
			err = utilError.FunctionError(attributeParams.Name + " - " + v2.ErrNoResult.Error())
		case 1:
			attributeId = response.Data[0].Id
		default:
			err = utilError.FunctionError(attributeParams.Name + " - " + v2.ErrTooManyResults.Error())
		}

		return attributeId, err

	} else {

		return existingAttributeId, err
	}

	return
}

type getRoleIdValue struct {
	Data []struct {
		RoleId int `json:"id,string"`
	} `json:"data"`
}

func (c *Client) GetRoleIdByName(roleName string) (roleId int, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": roleName,
	}

	response := getRoleIdValue{}
	err = c.v2.Query("int_getRoleIdByRoleName", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return 0, err
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(roleName + " - " + v2.ErrNoResult.Error())
	case 1:
		roleId = response.Data[0].RoleId
	default:
		err = utilError.FunctionError(roleName + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

func (c *Client) SetAttributeRole(attributeName string, roleName string, permission string) (err error) {

	var attributeID int
	var roleID int
	var permissionRead int
	var permissionWrite int

	attributeID, err = c.GetAttributeIdByAttributeName(attributeName)
	if err != nil {
		return err
	}

	roleID, err = c.GetRoleIdByName(roleName)
	if err != nil {
		return err
	}

	switch permission {
	case "x":
		permissionRead = 0
		permissionWrite = 0
	case "r":
		permissionRead = 1
		permissionWrite = 0
	case "w":
		permissionRead = 1
		permissionWrite = 1
	case "r/w":
		permissionRead = 1
		permissionWrite = 1
	default:
		return errors.New("The entered permission string is wrong: " + permission + ". You must use x,r,w or r/w.")
	}

	params := map[string]string{
		"argv1": strconv.Itoa(attributeID),
		"argv2": strconv.Itoa(roleID),
		"argv3": strconv.Itoa(permissionRead),
		"argv4": strconv.Itoa(permissionWrite),
	}

	var resp interface{}
	err = c.v2.Query("int_setAttributeRole", &resp, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	return
}
