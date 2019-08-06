package maptile

import (
	"reflect"
	"strings"

	"github.com/go-courier/reflectx"
)

func StructToFields(v interface{}) map[string]FieldType {
	structType := reflectx.Deref(reflect.TypeOf(v))
	if structType.Kind() != reflect.Struct {
		return nil
	}
	fields := map[string]FieldType{}
	for i := 0; i < structType.NumField(); i++ {
		ft := structType.Field(i)
		name, ok := ft.Tag.Lookup("name")
		if ok {
			name = strings.SplitN(name, ",", 2)[0]
		}

		if name == "-" {
			continue
		}

		if name == "" {
			name = ft.Name
		}

		switch ft.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			fields[name] = FieldTypeNumber
		case reflect.String:
			fields[name] = FieldTypeString
		case reflect.Bool:
			fields[name] = FieldTypeBoolean
		}
	}

	return fields
}

func StructToProperties(v interface{}) map[string]interface{} {
	s := reflectx.Indirect(reflect.ValueOf(v))
	if s.Kind() != reflect.Struct {
		return nil
	}
	typ := s.Type()
	props := map[string]interface{}{}
	for i := 0; i < s.NumField(); i++ {
		ft := typ.Field(i)
		name, ok := ft.Tag.Lookup("name")
		omitempty := false
		if ok {
			omitempty = strings.Contains(name, "omitempty")
			name = strings.SplitN(name, ",", 2)[0]
		}

		if name == "-" {
			continue
		}

		if name == "" {
			name = ft.Name
		}

		v := s.Field(i).Interface()
		if omitempty && reflectx.IsEmptyValue(v) {
			continue
		}
		props[name] = v
	}
	return props
}
