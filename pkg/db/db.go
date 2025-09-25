package db

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/slog"
	"xorm.io/xorm"
	"xorm.io/xorm/log"

	"github.com/wzhanjun/go-echo-skeleton/pkg/config"
	gslog "github.com/wzhanjun/log-service/client"
)

var (
	engine *xorm.Engine
	once   sync.Once
)

func GetEngine() *xorm.Engine {
	once.Do(func() {
		if location, err := time.LoadLocation(config.Cfg.System.Location); err == nil {
			time.Local = location
		}

		cfg := config.Cfg.MySql
		// dns
		dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)
		// engine
		var err error
		engine, err = xorm.NewEngine("mysql", dns)
		if err != nil {
			panic(err)
		}
		engine.SetTZLocation(time.Local)
		engine.ShowSQL(true)
		if !config.Cfg.System.StartCron {
			engine.SetLogger(NewSlogAdapter())
		}
	})

	return engine
}

type slogAdapter struct {
}

func NewSlogAdapter() log.Logger {
	return &slogAdapter{}
}

func (s *slogAdapter) Debug(v ...interface{}) {
	gslog.Label("SQL").Debug(v...)
}

func (s *slogAdapter) Debugf(format string, v ...interface{}) {
	gslog.Label("SQL").Debugf(format, v...)
}

func (s *slogAdapter) Error(v ...interface{}) {
	gslog.Label("SQL").Error(v...)
}

func (s *slogAdapter) Errorf(format string, v ...interface{}) {
	gslog.Label("SQL").Errorf(format, v...)
}

func (s *slogAdapter) Info(v ...interface{}) {
	gslog.Label("SQL").Info(v...)
}

func (s *slogAdapter) Infof(format string, v ...interface{}) {
	gslog.Label("SQL").Infof(format, v...)
}

func (s *slogAdapter) Warn(v ...interface{}) {
	gslog.Label("SQL").Warn(v...)
}

func (s *slogAdapter) Warnf(format string, v ...interface{}) {
	gslog.Label("SQL").Warnf(format, v...)
}

func (s *slogAdapter) Level() log.LogLevel {
	switch gslog.Std().Level {
	case slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel:
		return log.LOG_ERR
	case slog.InfoLevel:
		return log.LOG_INFO
	case slog.WarnLevel:
		return log.LOG_WARNING
	case slog.DebugLevel, slog.TraceLevel:
		return log.LOG_DEBUG
	}
	return log.LOG_DEBUG
}

func (s *slogAdapter) SetLevel(l log.LogLevel) {
	level := slog.TraceLevel
	switch l {
	case log.LOG_DEBUG:
		level = slog.DebugLevel
	case log.LOG_INFO:
		level = slog.InfoLevel
	case log.LOG_WARNING:
		level = slog.WarnLevel
	case log.LOG_ERR:
		level = slog.ErrorLevel
	case log.LOG_OFF:
		level = slog.PanicLevel
	}
	slog.SetLogLevel(level)
}

func (s *slogAdapter) ShowSQL(show ...bool) {

}

func (s *slogAdapter) IsShowSQL() bool {
	return true
}
