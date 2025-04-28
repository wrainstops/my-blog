package common

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Mode string
		Port string
	}
	Mysql struct {
		Host     string
		Port     string
		Database string
		Username string
		Password string
		Charset  string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
}

var Conf *Config

// ReadConfig 读取配置文件
func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败：" + err.Error())
	}
	if err := v.Unmarshal(&Conf); err != nil {
		// 注：viper对于配置文件的键不区分大小写
		panic("配置文件反序列化失败：" + err.Error())
	}

	return Conf
}
