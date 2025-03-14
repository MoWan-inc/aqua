package gorm

import (
	"context"
	"github.com/MoWan-inc/aqua/pkg/api"
	"github.com/MoWan-inc/aqua/pkg/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
)

type Transaction interface {
	Commit() error
	RollBack() error
	Session() *gorm.DB
}

type DAO interface {
	// Begin 启动新事务
	Begin() Transaction
	// WithTransaction 使用制定事务操作
	WithTransaction(tx Transaction) DAO
	// Count 默认支持的操作
	Count(ctx context.Context, q *api.QueryRequest, opts ...OptionFunc) (int64, error)
	// List 查询操作
	List(ctx context.Context, q *api.QueryRequest, results any, opts ...OptionFunc) error
	Get(ctx context.Context, obj domain.Indexer, result any, opts ...OptionFunc) error
	// ListWithInClause 含有IN语句的查询操作，手动输入需要的column和in clause value进行查询
	ListWithInClause(ctx context.Context, results any, query string, inClause [][]any) error
	// Delete 级联删除说明：传入的对象需要将外键引用的属性也赋值，否则不会触发GORM级联删除
	Delete(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error
	Create(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error
	// Update 只更新非空字段
	Update(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error
	// Save 覆盖式更新
	Save(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error
}

type OptionFunc func(*gorm.DB) *gorm.DB

var JoinOption = func(models ...any) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		for _, model := range models {
			if tName, ok := model.(string); ok {
				db = db.Joins(tName)
				continue
			}
			// 一对一关系
			table, err := schema.Parse(model, &sync.Map{}, db.NamingStrategy)
			if err != nil {
				panic("parse table error, join failed")
			}
			db = db.Joins(table.Name)
		}
		return db
	}
}
