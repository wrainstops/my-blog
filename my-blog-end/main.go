package main

import (
	"io"
	"my_blog_back/common"
	_ "my_blog_back/docs"
	"my_blog_back/router"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()

	gin.SetMode(gin.DebugMode) // debugMode
	// gin.SetMode(gin.ReleaseMode) // releaseMode

	gin.DefaultWriter = io.MultiWriter(common.Global_LogFile, os.Stdout) // gin log
	common.InitSlogLogger()                                              // slog

	// 创建服务
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery()) // 日志、Recovery中间件

	// 连接数据库，调用common中的InitDB
	common.InitDB()

	// 访问地址，接收请求
	router.CollectRoute(r)

	// 使用viper修改端口
	port := viper.GetString("server.port")
	if port != "" {
		r.Run(":" + port)
	} else {
		r.Run()
	}
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
