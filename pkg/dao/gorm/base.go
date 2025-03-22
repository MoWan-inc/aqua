package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/MoWan-inc/aqua/pkg/api"
	"github.com/MoWan-inc/aqua/pkg/domain"
	"github.com/MoWan-inc/aqua/pkg/util/object"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

const (
	BaseDAOName = "BaseDAO"
)

var (
	NotExistsError = errors.New("not exists")
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
	result := b.conn.WithContext(ctx)
	for _, o := range opts {
		result = o(result)
	}
	result = result.Where(q.Query).Model(q.Query)
	// filter
	result = prepareFieldFilter(&q.Filter, result)
	if !object.IsEmpty(q.Not) {
		result.Not(q.Not)
	}
	result.Count(&count)
	if result.Error != nil {
		// todo 日志记录
		return 0, result.Error
	}
	return
}

func (b *BaseDAO) List(ctx context.Context, q *api.QueryRequest, results any, opts ...OptionFunc) error {
	result := b.conn.WithContext(ctx)
	opts = append(opts, GetOptions(q.Query)...)
	for _, o := range opts {
		result = o(result)
	}
	result = result.Where(q.Query)
	// filter, pagination, sorting
	result = prepareFieldFilter(&q.Filter, result)
	if !object.IsEmpty(q.Not) {
		result.Not(q.Not)
	}
	result = prepareLimit(&q.Pagination, result)
	result = prepareSorting(&q.Sorting, result)
	result.Find(results)
	if result.Error != nil {
		// todo 添加日志
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("base dao list %v error: %w", q, result.Error)
	}
	return nil
}

func (b *BaseDAO) Get(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	result := b.conn.WithContext(ctx)
	opts = append(opts, GetOptions(obj)...)
	for _, o := range opts {
		result = o(result)
	}
	result = result.Where(obj).First(obj)
	if result.Error != nil {
		// todo 日志
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return NotExistsError
		}
		return fmt.Errorf("base dao get %v error: %w", obj, result.Error)
	}
	return nil
}

func (b *BaseDAO) ListWithInClause(ctx context.Context, results any, query string, inClause [][]any) error {
	result := b.conn.WithContext(ctx)
	result = result.Where(query, inClause)
	result.Find(results)
	if result.Error != nil {
		// todo 日志
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("base dao list %v with in clause %v error: %w", query, inClause, result.Error)
	}
	return nil
}

func (b *BaseDAO) Delete(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	result := b.conn.WithContext(ctx)
	for _, o := range opts {
		result = o(result)
	}
	// 如果有级联删除的对象，一起删除
	if result.Error != nil {
		// todo 日志
		return result.Error
	}
	return nil
}

func (b *BaseDAO) Create(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	result := b.conn.WithContext(ctx)
	for _, o := range opts {
		result = o(result)
	}
	result = b.conn.Create(obj)
	if result.Error != nil {
		// todo 日志
		return result.Error
	}
	return nil
}

func (b *BaseDAO) Update(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	// 增量覆盖更新，先找寻更新对象
	if err := b.updateByIndexer(obj); err != nil {
		return err
	}
	result := b.conn.WithContext(ctx)
	for _, o := range opts {
		result = o(result)
	}

	key := reflect.ValueOf(obj.Key())
	if key.IsZero() {
		return NotExistsError
	}
	result = result.Updates(obj)
	if result.Error != nil {
		// todo 日志
		return result.Error
	}
	return nil
}

func (b *BaseDAO) updateByIndexer(q domain.Indexer) error {
	// 主键存在则不需要找寻
	if !reflect.ValueOf(q.Key()).IsZero() {
		return nil
	}
	uniq, ok := q.(domain.UniqueIndex)
	if !ok {
		return nil
	}
	indexer := uniq.UniqIndexer()
	if indexer == nil {
		return nil
	}
	result := b.conn.Unscoped().Where(indexer).First(indexer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if result.Error != nil {
		// todo 日志
		return result.Error
	}
	if !reflect.ValueOf(indexer.Key()).IsZero() {
		// 存在单一索引，就更新该值
		if err := q.SetKey(indexer.Key()); err != nil {
			// todo 日志
			return err
		}
	}
	return nil
}

func (b *BaseDAO) Save(ctx context.Context, obj domain.Indexer, opts ...OptionFunc) error {
	// save表示全量更新，先找寻更新对象，没有找到则创建，这里由gorm的save实现
	if err := b.updateByIndexer(obj); err != nil {
		return err
	}
	result := b.conn.WithContext(ctx)
	for _, o := range opts {
		result = o(result)
	}
	result = b.conn.Save(obj)
	if result.Error != nil {
		// todo 日志
		return result.Error
	}
	return nil
}

func prepareFieldFilter(filter *api.Filter, result *gorm.DB) *gorm.DB {
	if len(filter.Filters) > 0 && len(filter.Filters) > 0 {
		filter.Filters = fmt.Sprintf("%%%v%%", filter.Filters)
		clauses := make([]string, 0)
		params := make([]any, 0)
		for _, f := range strings.Split(filter.Fields, ",") {
			clauses = append(clauses, fmt.Sprintf("%v LIKE ?", f))
			params = append(params, filter.Filters)
		}
		result = result.Where(strings.Join(clauses, " OR "), params...)
	}
	return result
}

func prepareLimit(pagination *api.Pagination, result *gorm.DB) *gorm.DB {
	if pagination.PageSize > 0 {
		result = result.Limit(pagination.PageSize).Offset((pagination.Page - 1) * pagination.PageSize)
	}
	return result
}

func prepareSorting(sorting *api.Sorting, result *gorm.DB) *gorm.DB {
	if len(sorting.SortBy) > 0 {
		order := sorting.SortBy
		if sorting.SortDesc {
			order += " DESC"
		}
		result = result.Order(order)
	}
	return result
}
