package repo

import (
	"github.com/wzhanjun/go-echo-skeleton/internal/dto"
	"xorm.io/xorm"
)

// QueryBuilder 对 session 施加 WHERE / ORDER / JOIN 等条件。
// 用于 Find / Paginate，避免为每个模型创建子 repo。
type QueryBuilder func(*xorm.Session) *xorm.Session

type BaseRepo[T any] interface {
	Create(session *xorm.Session, data *T) error
	Inserts(session *xorm.Session, data interface{}) error
	GetById(session *xorm.Session, id interface{}) (*T, error)
	Delete(session *xorm.Session, id interface{}) error
	Update(session *xorm.Session, id interface{}, data *T) error
	UpdateCols(session *xorm.Session, id interface{}, data *T, cols []string) error
	All(session *xorm.Session, params map[string]interface{}) ([]*T, error)
	AllWithBuilder(session *xorm.Session, builder QueryBuilder) ([]*T, error)
	Find(session *xorm.Session, builder QueryBuilder) ([]*T, error)
	GetBy(session *xorm.Session, builder QueryBuilder) (*T, error)
	Count(session *xorm.Session, builder QueryBuilder) (int64, error)
	Paginate(session *xorm.Session, page dto.PageParams, builder QueryBuilder) ([]*T, int64, error)
	FetchWithCache(session *xorm.Session, key string, builder QueryBuilder) ([]*T, error)
}
