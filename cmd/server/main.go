package main

import (
	"log"

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

	// config
	_ = config.Viper(*cfgDir, "yaml")
	// router
	e := routers.Router()

	// cron
	if config.Cfg.System.StartCron {
		if err := jobs.StartJob(); err != nil {
			log.Printf("start cron job failed, err:%+v", err)
			return
		}
	}

	// 启动地址
	e.Logger.Info("服务启动，地址：" + config.Cfg.System.Addr)
	e.Logger.Fatal(e.Start(config.Cfg.System.Addr))
}
