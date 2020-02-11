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
		case "[]int":
			err = bindAttrToIntSliceField(attr, valueField)
			if err != nil {
				return
			}
		default:
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
			SrcValue:  attr.Value,
		}
	}

	return
}

func bindAttrToIntSliceField(attr CiAttribute, field reflect.Value) (err error) {
	var numbers []int
	values := strings.Split(attr.Value, ",")

	for _, value := range values {
		trimmedValue := strings.TrimSpace(value)
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
