package api

import (
	"fmt"
	"reflect"
)

type BaseResponse[T any] struct {
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type internalList[T any] struct {
	Total int `json:"total"`
	List  []T `json:"list"`
}

type BaseListResponse[T any] struct {
	Msg  string          `json:"msg"`
	Data internalList[T] `json:"data"`
}

type ListResponse interface {
	New() ListResponse
	Merge(response ListResponse) error
	Length() int
	Total() int
}

func (rsp *BaseListResponse[T]) Total() int {
	return rsp.Data.Total
}
func (rsp *BaseListResponse[T]) Length() int {
	return len(rsp.Data.List)
}

func (rsp *BaseListResponse[T]) New() ListResponse {
	return &BaseListResponse[T]{
		Msg: "",
		Data: internalList[T]{
			List: []T{},
		},
	}
}

func (rsp *BaseListResponse[T]) Merge(other ListResponse) error {
	o, ok := other.(*BaseListResponse[T])
	if !ok {
		return fmt.Errorf("response merge error: %v", other)
	}
	rsp.Msg = o.Msg
	rsp.Data.List = append(rsp.Data.List, o.Data.List...)
	rsp.Data.Total += o.Data.Total
	return nil
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

// Get
type TemplateResponse BaseListResponse[Template]

// List
type TemplateListResponse BaseListResponse[Template]
