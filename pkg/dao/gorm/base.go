package gorm

import (
	"context"
	"github.com/MoWan-inc/aqua/pkg/api"
	"github.com/MoWan-inc/aqua/pkg/domain"
	"gorm.io/gorm"
)

const (
	BaseDAOName = "BaseDAO"
)

var _ DAO = &BaseDAO{}

type BaseDAO struct {
	conn *gorm.DB
}

func NewBaseDAO(db *gorm.DB) *BaseDAO {
	return &BaseDAO{conn: db}
}

func (b *BaseDAO) Name() string {
	return BaseDAOName
}

func (b *BaseDAO) Begin() Transaction {
	return &BaseDAO{conn: b.conn.Begin()}
}

func (b *BaseDAO) WithTransaction(tx Transaction) DAO {
	return &BaseDAO{conn: tx.Session()}
}

func (b *BaseDAO) Commit() error {
	return b.conn.Commit().Error
}

func (b *BaseDAO) RollBack() error {
	return b.conn.Rollback().Error
}

func (b *BaseDAO) Session() *gorm.DB {
	return b.conn
}

func (b *BaseDAO) Count(ctx context.Context, q *api.QueryRequest, opts ...OptionFunc) (count int64, err error) {
	return
}

func (b *BaseDAO) List(ctx context.Context, q *api.QueryRequest, results any, opts ...OptionFunc) error {
	return nil
}

func (b *BaseDAO) Get(ctx context.Context, obj domain.Indexer, result any, opts ...OptionFunc) error {
	return nil
}

func (b *BaseDAO) ListWithInClause(ctx context.Context, results any, query string, inClause [][]any) error {
	return nil
}

func (b *BaseDAO) Delete(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	return nil
}

func (b *BaseDAO) Create(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	return nil
}

func (b *BaseDAO) Update(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	return nil
}

func (b *BaseDAO) Save(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	return nil
}
