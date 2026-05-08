package services

import (
	"github.com/wzhanjun/go-echo-skeleton/internal/dto"
	"github.com/wzhanjun/go-echo-skeleton/internal/models"
	"github.com/wzhanjun/go-echo-skeleton/internal/repo/impl"
	"github.com/wzhanjun/go-echo-skeleton/pkg/db"
	"xorm.io/xorm"
)

type UserService struct {
	repo *impl.UserRepo
}

func NewUserService() *UserService {
	return &UserService{repo: impl.NewUserRepo()}
}

func (s *UserService) Get(id interface{}) (*models.User, error) {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.GetById(session, id)
}

func (s *UserService) FindByUsername(name string) (*models.User, error) {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.FindByUsername(session, name)
}

func (s *UserService) List(params dto.PageParams) ([]*models.User, int64, error) {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Paginate(session, params, func(sess *xorm.Session) *xorm.Session {
		return sess.Where("status = ?", 1).OrderBy("id DESC")
	})
}

func (s *UserService) Create(user *models.User) error {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Create(session, user)
}

func (s *UserService) Update(id interface{}, user *models.User) error {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Update(session, id, user)
}

func (s *UserService) Delete(id interface{}) error {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Delete(session, id)
}
