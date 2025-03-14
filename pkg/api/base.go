package api

import (
	"github.com/MoWan-inc/aqua/pkg/domain"
	"github.com/MoWan-inc/aqua/pkg/util/object"
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
	panic("implement me")
}
