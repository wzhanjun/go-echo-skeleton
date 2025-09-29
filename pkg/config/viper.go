package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper(path string, configType string) *viper.Viper {
	v := viper.New()
	// 设置文件

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

	// 设置类型
	v.SetConfigType(configType)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Fatal error config path: %s %v \n", path, err))
	}
	// 监听文件变化
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

	return v
}
