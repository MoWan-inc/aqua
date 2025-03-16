package api

import (
	"fmt"
	"github.com/MoWan-inc/aqua/pkg/domain"
	"github.com/MoWan-inc/aqua/pkg/util/object"
	"strings"
)

type QueryRequest struct {
	Query any `json:"query"`
	Pagination
	Sorting
	Filter
	Not any `json:"not"`
}

func QueryFor[T any](model T) *QueryRequest {
	return &QueryRequest{
		Query: model,
		Not:   object.NewObject[T](),
	}
}

type SaveRequest struct {
	domain.Indexer
}

type DeleteRequest struct {
	domain.Indexer
}

func SaveFor(model domain.Indexer) *SaveRequest {
	return &SaveRequest{Indexer: model}
}

func DeleteFor(model domain.Indexer) *DeleteRequest {
	return &DeleteRequest{Indexer: model}
}

func (r *QueryRequest) Validate() error {
	var fields map[string]any
	if len(r.SortBy) > 0 || len(r.Fields) > 0 {
		fields = domain.GetGormFields(object.ClassName(r.Query))
	}
	// sorting
	if len(r.SortBy) > 0 {
		if _, ok := fields[strings.TrimSpace(r.SortBy)]; !ok {
			return fmt.Errorf("query option error, invalid sort_by in fields: %s", r.SortBy)
		}
	}
	// filter
	if len(r.Fields) > 0 {
		filters := strings.Split(r.Fields, ",")
		for _, filter := range filters {
			filter = strings.TrimSpace(filter)
			//
			if strings.Contains(filter, ".") {
				filterTks := strings.Split(filter, ".")
				if len(filterTks) > 0 {
					filter = filterTks[len(filterTks)-1]
				}
			}
			if _, ok := fields[filter]; !ok {
				return fmt.Errorf("query option error, invalid filter %s in fields: %v", filter, r.Fields)
			}
		}
	}
	return nil
}

func (r *SaveRequest) Validate() error {
	// 数据合法
	if v, ok := r.Indexer.(Validator); ok {
		return v.Validate()
	}
	return nil
}

func (r *DeleteRequest) Validate() error {
	// 不允许不带任何条件删除全表
	if object.IsEmpty(r.Indexer) {
		return fmt.Errorf("delete condition error, empty condition")
	}
	return nil
}

type Validator interface {
	Validate() error
}
