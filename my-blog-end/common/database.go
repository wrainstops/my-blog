package common

import (
	"context"
	"fmt"
	"my_blog_back/model"
	"os/exec"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func InitDB(conf *Config) *gorm.DB {
	// 从反序列化后的配置文件中获取数据库配置
	c := conf.Mysql
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset)
	// 需事先引入mysql
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 数据迁移时不生成外键
	})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	// 自动创建表
	err2 := db.AutoMigrate(&model.User{}, &model.Article{}, &model.Like{})
	if err2 != nil {
		return nil
	}

	// DB实例赋值，即连接数据库
	DB = db
	return db
}

// GetDB 获取DB实例
func GetDB() *gorm.DB {
	return DB
}

func InitRedis(conf *Config) *redis.Client {
	cmd := exec.Command("redis-server")
	err := cmd.Start()
	if err != nil {
		panic("failed to start redis-server" + err.Error())
	}
	time.Sleep(time.Second * 3)

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,     // Redis 服务地址和端口
		Password: conf.Redis.Password, // Redis 密码
		DB:       conf.Redis.DB,       // 选择数据库
	})
	res, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("failed to connect redis" + err.Error())
	}
	fmt.Println("redis", rdb, "\n", res)

	RDB = rdb
	return rdb
}

func GetRedis() *redis.Client {
	return RDB
}
