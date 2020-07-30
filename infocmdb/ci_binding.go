package infocmdb

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
)

type BindError struct {
	Msg       string
	FieldName string
	SrcName   string
	SrcType   string
	SrcValue  string
}

func (e *BindError) Error() string {
	return e.Msg
}

func (c *Client) GetAndBindCi(ciId int, out interface{}) (err error) {
	attributes, err := c.GetCiAttributes(ciId)
	if err != nil {
		return
	}

	return bindCi(ciId, attributes, out)
}

func (c *Client) GetAndBindListOfCis(ciIds []int, out interface{}) (err error) {
	ciIdToAttributesMap, err := c.GetMapOfCiAttributes(ciIds)
	if err != nil {
		return
	}

	outSlicePtr := reflect.ValueOf(out)
	if outSlicePtr.Kind() != reflect.Ptr {
		return errors.New("out parameter is not a slice pointer")
	}
	outSlice := outSlicePtr.Elem()
	if outSlice.Kind() != reflect.Slice {
		return errors.New("out parameter is not a slice pointer")
	}
	outSliceElem := reflect.TypeOf(outSlice.Interface()).Elem()
	outSliceValue := reflect.MakeSlice(outSlice.Type(), 0, 0)

	for _, ciId := range ciIds {
		ciAttributes := ciIdToAttributesMap[ciId]

		elem := reflect.New(outSliceElem)
		err = bindCi(ciId, ciAttributes, elem.Interface())
		if err != nil {
			return
		}

		outSliceValue = reflect.Append(outSliceValue, elem.Elem())
	}

	outSlice.Set(outSliceValue)
	return
}

func (c *Client) GetAndBindListOfCisOfCiTypeName(ciTypeName string, out interface{}) (err error) {
	ciIds, err := c.GetListOfCiIdsOfCiTypeName(ciTypeName)
	if err != nil {
		err = errors.New("failed to get \"" + ciTypeName + "\" ci ids: " + err.Error())
		return
	}

	err = c.GetAndBindListOfCis(ciIds, out)
	return
}

func (c *Client) GetAndBindCiByAttributeValue(name string, value string, valueType v2.AttributeValueType, out interface{}) (err error) {
	ciId, err := c.GetCiIdByAttributeValue(name, value, valueType)
	if err != nil {
		return
	}

	err = c.GetAndBindCi(ciId, out)
	if err != nil {
		return
	}

	return
}

func (c *Client) GetAndBindCiByAttributeValueText(name string, value string, out interface{}) (err error) {
	return c.GetAndBindCiByAttributeValue(name, value, v2.ATTRIBUTE_VALUE_TYPE_TEXT, out)
}

func (c *Client) GetAndBindCiByAttributeValueCi(name string, value string, out interface{}) (err error) {
	return c.GetAndBindCiByAttributeValue(name, value, v2.ATTRIBUTE_VALUE_TYPE_CI, out)
}

func bindCi(ciId int, attributes []CiAttribute, out interface{}) (err error) {
	attrNameToAttrMap := map[string][]CiAttribute{}
	for _, attr := range attributes {
		ciAttributes := attrNameToAttrMap[attr.AttributeName]
		ciAttributes = append(ciAttributes, attr)
		attrNameToAttrMap[attr.AttributeName] = ciAttributes
	}

	outValue := reflect.ValueOf(out)
	for outValue.Kind() == reflect.Ptr || outValue.Kind() == reflect.Interface {
		outValue = outValue.Elem()
	}
	for i := 0; i < outValue.NumField(); i++ {
		structField := outValue.Type().Field(i)
		valueField := outValue.Field(i)

		ciTag := structField.Tag.Get("ci")
		if ciTag == "id" {
			valueField.SetInt(int64(ciId))
			continue
		}

		attrTag := structField.Tag.Get("attr")
		if attrTag == "" || attrTag == "-" {
			continue
		}

		attrs := attrNameToAttrMap[attrTag]

		err = bindAttr(attrs, structField, valueField)
		if err != nil {
			return
		}
	}

	return
}

func bindAttr(attrs []CiAttribute, structField reflect.StructField, valueField reflect.Value) (err error) {
	structFieldTypeName := structField.Type.String()

	switch structFieldTypeName {
	case "string":
		err = bindAttrToStringField(attrs, valueField)
		if err != nil {
			return
		}
	case "int":
		err = bindAttrToIntField(attrs, valueField)
		if err != nil {
			return
		}
	case "[]string":
		err = bindAttrToStringSliceField(attrs, valueField)
		if err != nil {
			return
		}
	case "[]int":
		err = bindAttrToIntSliceField(attrs, valueField)
		if err != nil {
			return
		}
	default:
		var attr CiAttribute
		if len(attrs) > 0 {
			attr = attrs[0]
		}
		return &BindError{
			Msg: fmt.Sprintf("failed to map struct field %v of type %v",
				structField.Name, structFieldTypeName),
			FieldName: structField.Name,
			SrcName:   attr.AttributeName,
			SrcType:   attr.AttributeType,
			SrcValue:  attr.Value,
		}
	}

	return
}

func bindAttrToStringField(attrs []CiAttribute, field reflect.Value) (err error) {
	if len(attrs) == 0 {
		return
	} else if len(attrs) == 1 {
		field.SetString(attrs[0].Value)
	} else {
		attr := attrs[0]
		return &BindError{
			Msg: fmt.Sprintf("failed to map multiple attributes with name \"%v\" to string",
				attr.AttributeName),
			FieldName: field.Type().Name(),
			SrcName:   attr.AttributeName,
			SrcType:   attr.AttributeType,
		}
	}

	return
}

func bindAttrToIntField(attrs []CiAttribute, field reflect.Value) (err error) {
	if len(attrs) == 0 {
		return
	} else if len(attrs) == 1 {
		attr := attrs[0]
		intValue, err := strconv.Atoi(attr.Value)
		if err != nil {
			return &BindError{
				Msg: fmt.Sprintf("failed to map attribute with name \"%v\" and value \"%v\" to int",
					attr.AttributeName, attr.Value),
				FieldName: field.Type().Name(),
				SrcName:   attr.AttributeName,
				SrcType:   attr.AttributeType,
			}
		}
		field.SetInt(int64(intValue))
	} else {
		attr := attrs[0]
		return &BindError{
			Msg: fmt.Sprintf("failed to map multiple attributes with name \"%v\" to int",
				attr.AttributeName),
			FieldName: field.Type().Name(),
			SrcName:   attr.AttributeName,
			SrcType:   attr.AttributeType,
		}
	}

	return
}

func bindAttrToStringSliceField(attrs []CiAttribute, field reflect.Value) (err error) {
	if len(attrs) == 0 {
		return
	}

	values := field.Interface().([]string)

	for _, attr := range attrs {
		switch attr.AttributeType {
		case "input":
			values = append(values, strings.TrimSpace(attr.Value))
		case "textarea":
			textareaValues := strings.Split(attr.Value, "\n")
			for _, textareaValue := range textareaValues {
				values = append(values, strings.TrimSpace(textareaValue))
			}
		default:
			return &BindError{
				Msg: fmt.Sprintf("failed to map attribute type %v to []string",
					attr.AttributeType),
				FieldName: field.Type().Name(),
				SrcName:   attr.AttributeName,
				SrcType:   attr.AttributeType,
				SrcValue:  attr.Value,
			}
		}
	}

	field.Set(reflect.ValueOf(values))
	return
}

func bindAttrToIntSliceField(attrs []CiAttribute, field reflect.Value) (err error) {
	if len(attrs) == 0 {
		return
	}

	// TODO: binding multiple attributes to int slice not implemented yet
	attr := attrs[0]

	var numbers []int
	values := strings.Split(attr.Value, ",")

	for _, value := range values {
		trimmedValue := strings.TrimSpace(value)
		if trimmedValue == "" {
			continue
		}

		number, err := strconv.Atoi(trimmedValue)
		if err != nil {
			return &BindError{
				Msg: fmt.Sprintf("failed convert attribute value \"%v\" to []int: %v",
					trimmedValue, err.Error()),
				FieldName: field.Type().Name(),
				SrcName:   attr.AttributeName,
				SrcType:   attr.AttributeType,
				SrcValue:  attr.Value,
			}
		}

		numbers = append(numbers, number)
	}

	field.Set(reflect.ValueOf(numbers))
	return
}
