package main

import (
	"io"
	"my_blog_back/common"
	_ "my_blog_back/docs"
	"my_blog_back/router"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 读取配置文件
	c := common.ReadConfig("./application.yml")

	// 配置gin模式
	gin.SetMode(c.Server.Mode)

	// gin log
	gin.DefaultWriter = io.MultiWriter(common.Global_LogFile, os.Stdout)
	// slog
	common.InitSlogLogger()

	// 创建服务
	r := gin.New()
	// 日志、Recovery中间件
	r.Use(gin.Logger(), gin.Recovery())

	// 连接数据库
	db := common.InitDB(c)
	rdb := common.InitRedis(c)

	go common.SyncRdbToMysql(db, rdb)

	// 访问地址，接收请求
	router.CollectRoute(r)

	// 端口，没有则默认8080
	port := c.Server.Port
	if port != "" {
		r.Run(":" + port)
	} else {
		r.Run()
	}
}
