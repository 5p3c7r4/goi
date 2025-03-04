package goi

import (
	"fmt"
	"reflect"
	"strconv"
)

func CustomDecode(src map[string]any, dst any) error {
	dstVal := reflect.ValueOf(dst).Elem()
	dstType := reflect.TypeOf(dst).Elem()

	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Field(i)
		fieldType := dstType.Field(i)
		tag := fieldType.Tag.Get("goi")

		if value, exists := src[tag]; exists {
			switch field.Kind() {
			case reflect.Pointer:
				if fieldType.Type.Elem() != reflect.TypeOf(value) {
					panic(fmt.Sprintf("bad type, source type %s destination type %s", reflect.TypeOf(value).String(), fieldType.Type.Elem().String()))
				}
				castedValue := reflect.ValueOf(value).Convert(reflect.TypeOf(value))
				ptr := reflect.New(reflect.TypeOf(value)).Elem()
				ptr.Set(castedValue)
				field.Set(reflect.ValueOf(ptr.Addr().Interface()))
			case reflect.String:
				if strVal, ok := value.(string); ok {
					field.SetString(strVal)
				} else {
					panic(fmt.Sprintf("expected string for %s, got %T", tag, value))
				}
			case reflect.Int:
				if intVal, ok := value.(int); ok {
					field.SetInt(int64(intVal))
				} else if strVal, ok := value.(string); ok {
					parsedInt, err := strconv.Atoi(strVal)
					if err != nil {
						panic(fmt.Sprintf("failed to convert string '%s' to int", strVal))
					}
					field.SetInt(int64(parsedInt))
				} else if floatVal, ok := value.(float64); ok {
					field.SetInt(int64(floatVal))
				} else {
					panic(fmt.Sprintf("expected int for %s, got %T", tag, value))
				}
			case reflect.Bool:
				if boolVal, ok := value.(bool); ok {
					field.SetBool(boolVal)
				} else if strVal, ok := value.(string); ok {
					parsedBool, err := strconv.ParseBool(strVal)
					if err != nil {
						panic(fmt.Sprintf("failed to convert string '%s' to bool", strVal))
					}
					field.SetBool(parsedBool)
				} else {
					panic(fmt.Sprintf("expected bool for %s, got %T", tag, value))
				}
			default:
				panic(fmt.Sprintf("unsupported type %s for field %s", field.Kind(), fieldType.Name))
			}
		}
	}

	return nil
}
