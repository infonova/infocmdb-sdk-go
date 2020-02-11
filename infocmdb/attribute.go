package infocmdb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"
	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

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

type BindError struct {
	Msg       string
	FieldName string
	SrcName   string
	SrcType   string
}

func (e *BindError) Error() string {
	return e.Msg
}

func (c *Client) GetAndBindCiAttributes(ciID int, out interface{}) (err error) {
	attributes, err := c.GetCiAttributes(ciID)
	if err != nil {
		return
	}

	return bindCiAttributes(attributes, out)
}

func bindCiAttributes(attributes []CiAttribute, out interface{}) (err error) {
	attrNameToAttrMap := map[string]CiAttribute{}
	for _, attr := range attributes {
		attrNameToAttrMap[attr.AttributeName] = attr
	}

	outValue := reflect.ValueOf(out)
	for outValue.Kind() == reflect.Ptr || outValue.Kind() == reflect.Interface {
		outValue = outValue.Elem()
	}
	for i := 0; i < outValue.NumField(); i++ {
		structField := outValue.Type().Field(i)
		valueField := outValue.Field(i)

		tag := structField.Tag.Get("attr")

		if tag == "" || tag == "-" {
			continue
		}

		attr := attrNameToAttrMap[tag]

		structFieldTypeName := structField.Type.String()
		switch structFieldTypeName {
		case "string":
			valueField.SetString(attr.Value)
		case "[]string":
			err = bindAttrToStringSliceField(attr, valueField)
			if err != nil {
				return
			}
		default:
			return &BindError{
				Msg:       fmt.Sprintf("failed to map struct field %v of type %v", structField.Name, structFieldTypeName),
				FieldName: structField.Name,
				SrcName:   attr.AttributeName,
				SrcType:   attr.AttributeType,
			}
		}
	}

	return
}

func bindAttrToStringSliceField(attr CiAttribute, field reflect.Value) (err error) {
	switch attr.AttributeType {
	case "textarea":
		values := strings.Split(attr.Value, "\n")

		for i, value := range values {
			values[i] = strings.TrimSpace(value)
		}

		field.Set(reflect.ValueOf(values))
	default:
		return &BindError{
			Msg:       fmt.Sprintf("failed to map attribute type %v to []string", attr.AttributeType),
			FieldName: field.Type().Name(),
			SrcName:   attr.AttributeName,
			SrcType:   attr.AttributeType,
		}
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

type CreateAttribute struct {
	CiID int `json:"id"`
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

func (c *Client) UpdateCiAttribute(ci int, ua []v2.UpdateCiAttribute) (err error) {
	return c.v2.UpdateCiAttribute(ci, ua)
}
