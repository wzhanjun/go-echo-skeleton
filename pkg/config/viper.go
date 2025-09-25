package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper(path string, configType string) *viper.Viper {
	v := viper.New()
	// 设置文件
	v.AddConfigPath("config")
	v.AddConfigPath("../config")
	v.AddConfigPath(".")

	if path != "" {
		v.SetConfigFile(path)
	}
	// 设置类型
	v.SetConfigType(configType)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Fatal error config file: %s %v \n", path, err))
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
