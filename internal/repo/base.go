package repo

import "xorm.io/xorm"

type BaseRepo[T any] interface {
	Create(session *xorm.Session, data *T) error
	Inserts(session *xorm.Session, data interface{}) error
	GetById(session *xorm.Session, id int) (*T, error)
	Update(session *xorm.Session, id int, data *T) error
	UpdateCols(session *xorm.Session, id int, data *T, cols []string) error
	All(session *xorm.Session, params map[string]interface{}) ([]*T, error)
	FetchWithCache(session *xorm.Session, key string, params map[string]interface{}) ([]*T, error)
}
