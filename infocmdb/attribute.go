package infocmdb

import (
	"strconv"

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
