package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init(path string, configType string) {
	// 1. 加载 .env (文件不存在则静默跳过)
	_ = godotenv.Load()

	v := viper.New()
	// 2. 允许环境变量覆盖配置 (例如 SYSTEM_ADDR 覆盖 system.addr)
	v.AutomaticEnv()

	if path != "" {
		info, err := os.Stat(path)
		if err != nil {
			panic(fmt.Errorf("config path error: %w", err))
		}
		if info.IsDir() {
			v.AddConfigPath(path)
		} else {
			v.SetConfigFile(path)
		}
	} else {
		v.AddConfigPath(".")
	}

	v.SetConfigType(configType)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Fatal error config path: %s %v \n", path, err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&Cfg); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&Cfg); err != nil {
		fmt.Println(err)
	}
}