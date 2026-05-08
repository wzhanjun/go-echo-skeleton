package models

type User struct {
	Id       int64  `xorm:"pk autoincr"`
	Username string `xorm:"varchar(50) notnull unique"`
	Status   int    `xorm:"default 1"`
}
