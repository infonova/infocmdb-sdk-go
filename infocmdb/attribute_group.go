package infocmdb

import (
	"strconv"
	"strings"

	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
)

var convertBoolToString = map[bool]string{false: "0", true: "1"}

type getAttributeGroupIdValue struct {
	Data []struct {
		GroupId int `json:"id,string"`
	} `json:"data"`
}

func (c *Client) GetAttributeGroupIdByName(attributeGroupName string) (attGroupId int, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetAttributeGroupIdByName_" + attributeGroupName
	cached, found := c.v1.Cache.Get(cacheKey)
	if found {
		return cached.(int), nil
	}

	params := map[string]string{
		"argv1": attributeGroupName,
	}

	response := getAttributeGroupIdValue{}
	err = c.v2.Query("int_getAttributeGroupIdByAttributeGroupName", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return 0, err
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(attributeGroupName + " - " + v2.ErrNoResult.Error())
	case 1:
		attGroupId = response.Data[0].GroupId
		c.v1.Cache.Set(cacheKey, attGroupId, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(attributeGroupName + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

type respCreateAttributeGroup struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []responseId `json:"data"`
}

type AttributeGroupParams struct {
	Name                     string
	Description              string
	Note                     string
	OrderNumber              int
	ParentAttributeGroupName string
	IsDuplicateAllow         bool
	IsActive                 bool
	UserId                   int
}

func (c *Client) NewAttributeGroupParams() (params *AttributeGroupParams) {
	params = &AttributeGroupParams{
		IsActive: true,
	}
	return
}

func (c *Client) CreateAttributeGroup(attributeGroupParams *AttributeGroupParams) (attributeGroupId int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	existingAttributeGroup, err := c.GetAttributeGroupIdByName(attributeGroupParams.Name)
	if err != nil && strings.Contains(err.Error(), "query returned no result") == false {
		return 0, err
	}

	if existingAttributeGroup == 0 {

		columns := []string{
			"name",
			"description",
			"note",
			"order_number",
			"parent_attribute_group_id",
			"is_duplicate_allow",
			"is_active",
			"user_id",
		}

		parentAttributeGroupId, err := c.GetAttributeGroupIdByName(attributeGroupParams.ParentAttributeGroupName)
		if err != nil {
			return 0, err
		}

		values := []string{
			attributeGroupParams.Name,
			attributeGroupParams.Description,
			attributeGroupParams.Note,
			strconv.Itoa(attributeGroupParams.OrderNumber),
			strconv.Itoa(parentAttributeGroupId),
			convertBoolToString[attributeGroupParams.IsDuplicateAllow],
			convertBoolToString[attributeGroupParams.IsActive],
			strconv.Itoa(attributeGroupParams.UserId),
		}

		params := map[string]string{
			"argv1": "`" + strings.Join(columns, "`, `") + "`",
			"argv2": "'" + strings.Join(values, "', '") + "'",
		}

		response := respCreateAttributeGroup{}
		err = c.v2.Query("int_createAttributeGroup", &response, params)
		if err != nil {
			err = utilError.FunctionError(err.Error())
			log.Error("Error: ", err)
			return 0, err
		}

		switch len(response.Data) {
		case 0:
			err = utilError.FunctionError(attributeGroupParams.Name + " - " + v2.ErrNoResult.Error())
		case 1:
			attributeGroupId = response.Data[0].Id
		default:
			err = utilError.FunctionError(attributeGroupParams.Name + " - " + v2.ErrTooManyResults.Error())
		}

		return attributeGroupId, err

	} else {

		return existingAttributeGroup, err
	}

	return
}
