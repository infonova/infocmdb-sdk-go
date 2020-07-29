package infocmdb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func bindCiAttributes(attributes []CiAttribute, out interface{}) (err error) {
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

		tag := structField.Tag.Get("attr")
		if tag == "" || tag == "-" {
			continue
		}

		attrs := attrNameToAttrMap[tag]

		structFieldTypeName := structField.Type.String()
		switch structFieldTypeName {
		case "string":
			err = bindAttrToStringField(attrs, valueField)
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
				Msg:       fmt.Sprintf("failed to map struct field %v of type %v", structField.Name, structFieldTypeName),
				FieldName: structField.Name,
				SrcName:   attr.AttributeName,
				SrcType:   attr.AttributeType,
				SrcValue:  attr.Value,
			}
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
			Msg:       fmt.Sprintf("failed to map multiple attributes with name \"%v\" to string", attr.AttributeName),
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
				Msg:       fmt.Sprintf("failed to map attribute type %v to []string", attr.AttributeType),
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
				Msg:       fmt.Sprintf("failed convert attribute value \"%v\" to []int: %v", trimmedValue, err.Error()),
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
