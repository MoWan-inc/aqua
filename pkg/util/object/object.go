package object

import "reflect"

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
