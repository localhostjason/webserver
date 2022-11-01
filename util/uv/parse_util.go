package uv

import (
	"reflect"
	"strings"
)

// Trim 使用方法在tag中添加 `trim:"left|right|ignore"`即可
func Trim(p interface{}) {
	if reflect.TypeOf(p).Kind() != reflect.Ptr {
		return
	}
	t := reflect.TypeOf(p).Elem()
	v := reflect.ValueOf(p).Elem()
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			trimNestedStruct(field)
			continue
		}
		if field.Kind() == reflect.String {
			field.Set(reflect.ValueOf(trim(field.String(), t.Field(i).Tag.Get("trim"))))
		}
	}
}

func trimNestedStruct(parentField reflect.Value) {
	structField := parentField.Type()
	for j := 0; j < structField.NumField(); j++ {
		if structField.Field(j).Type.Kind() == reflect.Struct {
			trimNestedStruct(parentField.Field(j))
			continue
		}
		if structField.Field(j).Type.Kind() == reflect.String {
			parentField.Field(j).Set(reflect.ValueOf(trim(parentField.Field(j).String(), structField.Field(j).Tag.Get("trim"))))
		}
	}
}

func trim(str string, trimType string) string {
	switch trimType {
	case "ignore":
		return str
	case "left":
		return strings.TrimLeft(str, " ")
	case "right":
		return strings.TrimRight(str, " ")
	default:
		return strings.Trim(str, " ")
	}
}
