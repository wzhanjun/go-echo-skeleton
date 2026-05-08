package impl

import (
	"context"

	"github.com/wzhanjun/go-echo-skeleton/internal/dto"
	"github.com/wzhanjun/go-echo-skeleton/internal/repo"
	"github.com/wzhanjun/go-echo-skeleton/pkg/cache"
	slog "github.com/wzhanjun/log-service/client"
	"xorm.io/xorm"
)

type baseRepo[T any] struct {
	cache cache.Cache
}

func NewBaseRepoImpl[T any]() repo.BaseRepo[T] {
	return NewBaseRepoWithCache[T](cache.NewRedisCache())
}

func NewBaseRepoWithCache[T any](c cache.Cache) repo.BaseRepo[T] {
	return &baseRepo[T]{cache: c}
}

func (s *baseRepo[T]) Create(session *xorm.Session, data *T) error {
	_, err := session.Insert(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Create failed, data:%+v, err:%+v", data, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) Inserts(session *xorm.Session, data interface{}) error {
	_, err := session.Insert(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Inserts failed, data:%+v, err:%+v", data, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) GetById(session *xorm.Session, id interface{}) (*T, error) {
	data := new(T)
	_, err := session.ID(id).Get(data)
	if err != nil {
		slog.Label("repo").Errorf("%T.GetById failed, id=%v, err:%+v", data, id, err)
		return nil, err
	}
	return data, nil
}

func (s *baseRepo[T]) Delete(session *xorm.Session, id interface{}) error {
	data := new(T)
	_, err := session.ID(id).Delete(data)
	if err != nil {
		slog.Label("repo").Errorf("%T.Delete failed, id=%v, err:%+v", data, id, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) Update(session *xorm.Session, id interface{}, data *T) error {
	_, err := session.ID(id).AllCols().Update(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Update failed, id:%v, data:%+v, err:%+v", id, data, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) UpdateCols(session *xorm.Session, id interface{}, data *T, cols []string) error {
	_, err := session.ID(id).Cols(cols...).Update(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.UpdateCols failed, id:%v, data:%+v, cols:%v, err:%+v", id, data, cols, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) All(session *xorm.Session, params map[string]interface{}) ([]*T, error) {
	list := make([]*T, 0)
	err := session.Where(params).Find(&list)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.All failed, params:%+v, err:%+v", params, err)
		return nil, err
	}
	return list, nil
}

func (s *baseRepo[T]) AllWithBuilder(session *xorm.Session, builder repo.QueryBuilder) ([]*T, error) {
	list := make([]*T, 0)
	err := builder(session).Find(&list)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.AllWithBuilder failed, err:%+v", err)
		return nil, err
	}
	return list, nil
}

func (s *baseRepo[T]) Find(session *xorm.Session, builder repo.QueryBuilder) ([]*T, error) {
	list := make([]*T, 0)
	err := builder(session).Find(&list)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Find failed, err:%+v", err)
		return nil, err
	}
	return list, nil
}

func (s *baseRepo[T]) GetBy(session *xorm.Session, builder repo.QueryBuilder) (*T, error) {
	data := new(T)
	_, err := builder(session).Get(data)
	if err != nil {
		slog.Label("repo").Errorf("%T.GetBy failed, err:%+v", data, err)
		return nil, err
	}
	return data, nil
}

func (s *baseRepo[T]) Count(session *xorm.Session, builder repo.QueryBuilder) (int64, error) {
	total, err := builder(session).Count(new(T))
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Count failed, err:%+v", err)
		return 0, err
	}
	return total, nil
}

func (s *baseRepo[T]) Paginate(session *xorm.Session, page dto.PageParams, builder repo.QueryBuilder) ([]*T, int64, error) {
	page.Normalize()

	list := make([]*T, 0)
	total, err := builder(session).
		Limit(page.PageSize, (page.PageIndex-1)*page.PageSize).
		FindAndCount(&list)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Paginate.Find failed, err:%+v", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (s *baseRepo[T]) FetchWithCache(session *xorm.Session, key string, builder repo.QueryBuilder) ([]*T, error) {
	load := func() ([]*T, error) {
		list := make([]*T, 0)
		if builder != nil {
			err := builder(session).Find(&list)
			return list, err
		}
		err := session.Find(&list)
		return list, err
	}

	return cache.Remember(context.Background(), s.cache, key, 300, load)
}
