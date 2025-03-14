package domain

import "time"

// Indexer 主键，通常自增ID
type Indexer interface {
	Key() any
	SerKey(key any) error
	UpdatedTime() time.Time
}

// UniqueIndex 唯一索引，实现这个接口表示数据库不会出现相同索引的数据
type UniqueIndex interface {
	UniqIndexer() Indexer
}
