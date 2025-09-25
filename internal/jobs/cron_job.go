package jobs

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/wzhanjun/go-echo-skeleton/pkg/config"
	slog "github.com/wzhanjun/log-service/client"
)

type CronJob interface {
	Serve()
	RegJob(exp string, fn func()) error
}

type jobContext struct {
	*cron.Cron
}

func (s *jobContext) Serve() {
	s.Cron.Start()
}

func (s *jobContext) RegJob(exp string, fn func()) error {
	id, err := s.AddFunc(exp, func() {
		defer func() {
			if p := recover(); p != nil {
				slog.Label("cron").Errorf("执行(%s)定时任务panic，%+v", exp, p)
			}
		}()
		fn()
	})
	if err != nil {
		slog.Label("cron").Errorf("注册定时任务(%s)失败，entryID:%d，err:%+v", exp, id, err)
		return err
	}
	return nil
}

func NewJob() CronJob {
	// location
	location, err := time.LoadLocation(config.Cfg.System.Location)
	if err != nil {
		location = time.Local
	}
	return &jobContext{
		Cron: cron.New(cron.WithSeconds(), cron.WithLocation(location)),
	}
}
