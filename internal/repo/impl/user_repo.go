package impl

import (
	"github.com/wzhanjun/go-echo-skeleton/internal/models"
	"github.com/wzhanjun/go-echo-skeleton/internal/repo"
	"xorm.io/xorm"
)

type UserRepo struct {
	repo.BaseRepo[models.User]
}

func (r *UserRepo) FindByUsername(session *xorm.Session, name string) (*models.User, error) {
	return r.GetBy(session, func(s *xorm.Session) *xorm.Session {
		return s.Where("username = ?", name)
	})
}

func NewUserRepo() *UserRepo {
	return &UserRepo{NewBaseRepoImpl[models.User]()}
}
