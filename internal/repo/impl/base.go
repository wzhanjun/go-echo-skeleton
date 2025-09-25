package impl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wzhanjun/go-echo-skeleton/internal/repo"
	"github.com/wzhanjun/go-echo-skeleton/internal/utils"
	"github.com/wzhanjun/go-echo-skeleton/pkg/cache"
	slog "github.com/wzhanjun/log-service/client"
	"xorm.io/xorm"
)

type baseRepo[T any] struct {
	cache *redis.Client
}

func NewBaseRepoImpl[T any]() repo.BaseRepo[T] {
	return &baseRepo[T]{
		cache: cache.GetRedisConn(),
	}
}

func (s *baseRepo[T]) Inserts(session *xorm.Session, data interface{}) error {
	_, err := session.Insert(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Inserts failed, data:%+v, err:%+v", data, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) Create(session *xorm.Session, data *T) error {
	_, err := session.Insert(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Create failed, data:%+v, err:%+v", data, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) GetById(session *xorm.Session, id int) (*T, error) {
	data := new(T)
	_, err := session.ID(id).Get(data)
	if err != nil {
		slog.Label("repo").Errorf("%T.Get failed, userId=%d, err:%+v", data, id, err)
		return nil, err
	}
	return data, nil
}

func (s *baseRepo[T]) Update(session *xorm.Session, id int, data *T) error {
	_, err := session.ID(id).AllCols().Update(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.Update failed, id:%d, data:%+v, err:%+v", id, data, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) UpdateCols(session *xorm.Session, id int, data *T, cols []string) error {
	_, err := session.ID(id).Cols(cols...).Update(data)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.UpdateCols failed, id:%d, data:%+v, cods:%v, err:%+v", id, data, cols, err)
		return err
	}
	return nil
}

func (s *baseRepo[T]) All(session *xorm.Session, params map[string]interface{}) ([]*T, error) {
	list := make([]*T, 0)
	err := session.Where(params).Find(&list)
	if err != nil {
		slog.Label("repo").Errorf("baseRepo.List query failed, list:%+v params:%+v,error:%+v", list, params, err)
		return nil, err
	}
	return list, err
}

func (s *baseRepo[T]) FetchWithCache(session *xorm.Session, key string, params map[string]interface{}) ([]*T, error) {

	list := make([]*T, 0)

	cVal := s.cache.Get(context.Background(), key).Val()
	if cVal == "" {
		err := session.Where(params).Find(&list)
		if err != nil {
			slog.Label("repo").Errorf("baseRepo.List query failed, list:%+v params:%+v,error:%+v", list, params, err)
			return nil, err
		}

		s.cache.Set(context.Background(), key, utils.ToJson(list), time.Minute*5)

		return list, nil
	}

	err := json.Unmarshal([]byte(cVal), &list)

	return list, err
}
