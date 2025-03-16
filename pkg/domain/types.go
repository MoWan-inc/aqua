package domain

import (
	"encoding/json"
	"fmt"
	"github.com/MoWan-inc/aqua/pkg/util/object"
	"go/ast"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

// Indexer 主键，通常自增ID
type Indexer interface {
	Key() any
	SetKey(key any) error
	UpdatedTime() time.Time
}

// UniqueIndex 唯一索引，实现这个接口表示数据库不会出现相同索引的数据
type UniqueIndex interface {
	UniqIndexer() Indexer
}

// 路由映射
var DomainPath = map[string]string{
	object.ClassName(Template{}): "template",
}

// Relation 关联关系，返回关联表名或对象，限于一对一关系
type Relation interface {
	Joins() []any
}

// Preload 预加载
// 可查看官方文档 https://gorm.io/zh_CN/docs/preload.html
// gorm中 has many关系即一对多关系的级联表需要preload加载级联表信息，官方参数是字段名字
type Preload interface {
	Preloads() []string
}

// Model 软删除模型
type Model struct {
	ID        uint           `form:"id" json:"id,omitempty" gorm:"column:id; primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
}

// 手写DEEP COPY
// 目前Go生成的代码不支持time.time
func (m *Model) DeepCopyInto(out *Model) {
	*out = *m
}

func (m *Model) DeepCopy() *Model {
	if m == nil {
		return nil
	}
	out := new(Model)
	m.DeepCopyInto(out)
	return out
}

func (m *Model) Key() any { return m.ID }

func (m *Model) SetKey(key any) error {
	switch key := key.(type) {
	case uint:
		m.ID = key
	case int:
		m.ID = uint(key)
	case *uint:
		m.ID = *key
	default:
		return fmt.Errorf("key type error, only uint supported: %v", key)
	}
	return nil
}

func (m *Model) UpdatedTime() time.Time { return m.UpdatedAt }

func (m *Model) Update(model Model) {
	if !model.CreatedAt.IsZero() && m.CreatedAt != model.CreatedAt {
		m.CreatedAt = model.CreatedAt
	}
	if !model.UpdatedAt.IsZero() && m.UpdatedAt != model.UpdatedAt {
		m.UpdatedAt = model.UpdatedAt
	}
	if model.DeletedAt.Valid && m.DeletedAt != model.DeletedAt {
		m.DeletedAt = model.DeletedAt
	}
}

// HardDeleteModel 硬删除模型
type HardDeleteModel struct {
	ID        uint      `form:"id" json:"id,omitempty" gorm:"column:id; primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// DeepCopyInto 深拷贝
func (m *HardDeleteModel) DeepCopyInto(out *HardDeleteModel) {
	*out = *m
}

func (m *HardDeleteModel) DeepCopy() *HardDeleteModel {
	if m == nil {
		return nil
	}
	out := new(HardDeleteModel)
	m.DeepCopyInto(out)
	return out
}

func (m *HardDeleteModel) Key() any {
	return m.ID
}

func (m *HardDeleteModel) SetKey(key any) error {
	switch key := key.(type) {
	case uint:
		m.ID = key
	case int:
		m.ID = uint(key)
	case *uint:
		m.ID = *key
	default:
		return fmt.Errorf("key type error, only uint supported: %v", key)
	}
	return nil
}

func (m *HardDeleteModel) UpdatedTime() time.Time { return m.UpdatedAt }

func (m *HardDeleteModel) Update(model HardDeleteModel) {
	if !model.CreatedAt.IsZero() && m.CreatedAt != model.CreatedAt {
		m.CreatedAt = model.CreatedAt
	}
	if !model.UpdatedAt.IsZero() && m.UpdatedAt != model.UpdatedAt {
		m.UpdatedAt = model.UpdatedAt
	}
}

/*
ColumnModel 列的基类，用于读写JSON对象。原因：数据库中存储json字符串，转换为go中的struct需要转换，这样就不用新增column就可以加字段
对于 json 列适用的注意事项
1. 列是对象的指针类型、数组时，直接使用gorm的序列化标签 serializer:json，无需继承这个基类
2. 列是对象类型时，需要继承该基类（因为默认的方法不支持null值转对象空值），同时还需要适用gorm的序列化标签 serializer:json
3. 只适用于列，不用于关联表
*/
type ColumnModel struct{}

func (c *ColumnModel) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch value := value.(type) {
	case []byte:
		return json.Unmarshal(value, c)
	case string:
		return json.Unmarshal([]byte(value), c)
	default:
		return fmt.Errorf("unsupported type: %+v", value)
	}
}

func (c *ColumnModel) Value() (any, error) {
	vl := reflect.ValueOf(c)
	if vl.IsZero() {
		return nil, nil
	}
	return json.Marshal(c)
}

const (
	DateLayout      = "20060102"
	GormTag         = "gorm"
	GormColumnKey   = "column"
	GormEmbeddedKey = "embedded"
)

// model 运行时不变，缓存起来即可
var gormModelFields = map[string]map[string]any{}
var gormModelTypes = map[string]reflect.Type{}

func init() {
	modelTypes := []reflect.Type{
		reflect.ValueOf(Template{}).Type(),
	}
	for _, modelType := range modelTypes {
		initGormFields(modelType)
		gormModelTypes[modelType.Name()] = modelType
	}
}

func getBaseFields() map[string]any {
	fields := map[string]any{}
	fields["id"] = struct{}{}
	fields["created_at"] = struct{}{}
	fields["updated_at"] = struct{}{}
	return fields
}

func initGormFields(modelType reflect.Type) map[string]any {
	fields := getBaseFields()
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		if !ast.IsExported(f.Name) {
			continue
		}
		gormTagDesc := f.Tag.Get(GormTag)
		if len(gormTagDesc) == 0 {
			continue
		}
		gormTags := strings.Split(gormTagDesc, ";")
		fieldColumn := ""
		embedFields := map[string]any{}
		for _, tag := range gormTags {
			if len(tag) == 0 {
				continue
			}
			tks := strings.Split(tag, ":")
			// 找到 column 定义
			if len(tks) == 2 && strings.TrimSpace(tks[0]) == GormColumnKey {
				fieldColumn = strings.TrimSpace(tks[1])
				continue
			}
			// 如果有 embedded 定义，递归查找嵌套属性
			if len(tks) == 1 && strings.TrimSpace(tks[0]) == GormEmbeddedKey {
				embedFields = initGormFields(f.Type)
				break
			}
		}

		if len(embedFields) > 0 {
			// 嵌套属性
			for key := range embedFields {
				fields[key] = struct{}{}
			}
		} else if len(fieldColumn) > 0 {
			// 无嵌套属性，则采用本列属性说明
			fields[fieldColumn] = struct{}{}
		} // 如果无列属性说明，则不作为gorm列
	}
	gormModelFields[modelType.Name()] = fields
	return fields
}

func GetGormFields(cls string) map[string]any {
	modelFields, ok := gormModelFields[cls]
	if ok {
		return modelFields
	}
	return map[string]any{}
}

func IsTypeValid(t string) bool {
	_, ok := gormModelTypes[t]
	return ok
}
