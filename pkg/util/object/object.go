package object

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// 更新对象的所有非空字段，浅拷贝
func AttributeUpdate(dest, src any) error {
	dType := reflect.TypeOf(dest)
	sType := reflect.TypeOf(src)
	if dType != sType {
		return fmt.Errorf("attribute update error, type mismatch: %v != %v", dType, sType)
	}

	dVal := reflect.ValueOf(dest).Elem()
	sVal := reflect.ValueOf(src).Elem()
	for i := dVal.NumField() - 1; i >= 0; i-- {
		if isZero(sVal.Field(i)) {
			continue
		}
		if !reflect.DeepEqual(sVal.Field(i).Interface(), dVal.Field(i).Interface()) {
			dVal.Field(i).Set(sVal.Field(i))
		}
	}
	return nil
}

// 判断新对象非空字段是否与老对象一致
func AttributeEqual(dest, src any) bool {
	dType := reflect.TypeOf(dest)
	sType := reflect.TypeOf(src)
	if dType != sType {
		return false
	}

	dVal := reflect.Indirect(reflect.ValueOf(dest))
	sVal := reflect.Indirect(reflect.ValueOf(src))
	for i := dVal.NumField() - 1; i >= 0; i-- {
		if isZero(sVal.Field(i)) {
			continue
		}
		if !reflect.DeepEqual(sVal.Field(i).Interface(), dVal.Field(i).Interface()) {
			return false
		}
	}
	return true
}

func isZero(obj reflect.Value) bool {
	if obj.IsZero() {
		return true
	}
	switch obj.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Chan:
		return obj.Len() == 0
	case reflect.Pointer, reflect.Struct:
		return isEmptyObject(obj)
	default:
		return false

	}
}

func IsEmpty(obj any) bool {
	return obj == nil || isZero(reflect.ValueOf(obj))
}

func isEmptyObject(dVal reflect.Value) bool {
	for dVal.Kind() == reflect.Pointer {
		dVal = dVal.Elem()
	}
	if dVal.Kind() != reflect.Struct {
		return false
	}
	// 递归判断结构体字段是否为空
	for i := dVal.NumField() - 1; i >= 0; i-- {
		if !isZero(dVal.Field(i)) {
			return false
		}
	}
	return true
}

/*
NewObject creates a new object of the given type.
例如： NewObject[Template]() = Template{}

	NewObject[*Template]() = &Template{}
*/
func NewObject[T any]() T {
	var obj T
	t := reflect.ValueOf(obj).Type()
	var n reflect.Value
	switch t.Kind() {
	case reflect.Pointer:
		n = reflect.New(t.Elem())
	case reflect.Struct:
		n = reflect.Zero(t)
	}
	v := reflect.Indirect(reflect.ValueOf(&obj))
	v.Set(n)
	return obj
}

func ClassName[T any](obj T) string {
	t := reflect.TypeOf(obj)
	for {
		k := t.Kind()
		switch k {
		case reflect.Map, reflect.Array, reflect.Pointer, reflect.Slice, reflect.Chan:
			t = t.Elem()
		default:
			return t.Name()
		}
	}
}

func MustMarshalJSON(obj any) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("marshal any to json error %v", err))
	}
	return string(bytes)
}
