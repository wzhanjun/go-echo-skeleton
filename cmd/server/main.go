package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/wzhanjun/go-echo-skeleton/internal/jobs"
	"github.com/wzhanjun/go-echo-skeleton/internal/routers"
	"github.com/wzhanjun/go-echo-skeleton/pkg/config"
)

var (
	cfgDir = pflag.StringP("config dir", "c", "../../config", "config path.")
)

func main() {
	pflag.Parse()

	config.Init(*cfgDir, "yaml")
	e := routers.Router()

	if config.Cfg.System.StartCron {
		if err := jobs.StartJob(); err != nil {
			log.Printf("start cron job failed, err:%+v", err)
			return
		}
	}

	// 启动 HTTP 服务 (goroutine 以便捕获关闭信号)
	go func() {
		e.Logger.Info("服务启动，地址：" + config.Cfg.System.Addr)
		if err := e.Start(config.Cfg.System.Addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	e.Logger.Info("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("服务已关闭")
}
