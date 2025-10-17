package model

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Password string `gorm:"size:255;not null"`

	Articles []Article
}

/**
 * AddUser 创建用户
 */
func AddUser(db *gorm.DB, user User) error {
	res := db.Create(&user)
	return res.Error
}

/**
 * GetUserById 用id搜索用户
 * @param userId {uint} 用户id
 * @return {User} 用户
 */
func GetUserById(db *gorm.DB, userId uint) (User, error) {
	var user User
	res := db.Where("id = ?", userId).First(&user)
	return user, res.Error
}

/**
 * GetUserByName 用name搜索用户
 * @param name {string} 用户name
 * @return {User} 用户
 */
func GetUserByName(db *gorm.DB, name string) (User, error) {
	var user User
	res := db.Where("name = ?", name).First(&user)
	return user, res.Error
}

/**
 * GetStat 用户统计信息
 * @param userId {uint} 用户id
 * @return {int64} 发布博客数量
 * @return {int64} 点赞量
 * @return {int64} 被回复量
 * @return {int64} 被点赞量
 */
func GetStat(db *gorm.DB, userId uint) (a, b, c, d int64, err error) {
	type Result1 struct {
		/*
			GORM在执行Scan操作时，要求结构体的字段名大写
			主要是因为GORM在处理查询结果时，默认会将查询到的列名映射到结构体的字段名。
			如果字段名小写，GORM无法正确映射这些字段，导致无法将查询结果填充到结构体中。
		*/
		C int64
		D int64
	}
	var result1 Result1

	res := db.Model(&Article{}).
		Select("COALESCE(SUM(reply_num), 0) AS be_replied_num", "COALESCE(SUM(like_num), 0) AS be_liked_num").
		Where("user_id = ?", userId).
		Count(&a).
		Scan(&result1)

	res = db.Model(&Like{}).
		Where("user_id = ?", userId).
		Count(&b)

	err = res.Error
	return
}

/**
 * GetRecentContacts 获取最近聊天的联系人
 * @param ctx {context.Context} redis context
 * @param key {string} redis key
 * @return {[]string} 最近聊天的联系人
 */
func GetRecentContacts(rdb *redis.Client, ctx context.Context, key string) []string {
	return rdb.LRange(ctx, key, 0, -1).Val()
}
