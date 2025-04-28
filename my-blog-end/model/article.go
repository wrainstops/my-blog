package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model

	/*
	 * gorm: gorm配置
	 * json: 查询到的数据转换成json格式后的label
	 */
	Title    string `gorm:"type:varchar(50);default:null;comment:博客标题" json:"title"`
	Content  string `gorm:"type:varchar(2000);not null;comment:博客内容/回复内容" json:"content"`
	ParentId uint   `gorm:"default:null;comment:父博客id(回复的父博客)" json:"parentId"`
	ReplyNum int64  `gorm:"comment:回复数量" json:"replyNum"`
	LikeNum  int64  `gorm:"comment:点赞数量" json:"likeNum"`

	User     User `gorm:"foreignKey:UserID"`
	UserID   uint `gorm:"not null;comment:作者id" json:"authId"` // 逻辑外键
	ToAuth   User `gorm:"foreignKey:ToAuthId"`
	ToAuthId uint `gorm:"default:null;comment:被回复者id" json:"toAuthId"` // 逻辑外键
}

/**
 * AddArticle 新增博客/回复
 * @param article {*Article} 博客/回复
 */
func AddArticle(db *gorm.DB, article *Article) error {
	res := db.Create(&article)
	return res.Error
}

/**
 * GetArticleOrReplyList 博客/回复列表(不需要parentId)
 * @param listType {uint} 0:博客列表 1:回复列表
 * @param key {string} 关键字
 * @param descCreatedTime {bool} true:降序 默认升序
 * @return {[]Article} 博客/回复列表
 * @return {int64} 统计总数
 */
func GetArticleOrReplyList(db *gorm.DB, page, pageSize int, listType uint, key string, descCreatedTime bool, userId uint) (list []Article, total int64, err error) {
	db = db.Model(&Article{})
	if key != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+key+"%", "%"+key+"%")
	}
	if descCreatedTime == true {
		db = db.Order("created_at desc")
	}
	if userId != 0 {
		db = db.Where("user_id = ?", userId)
	}
	if listType == 1 {
		db = db.Where("parent_id IS NOT NULL")
	} else {
		db = db.Where("parent_id IS NULL")
	}
	res := db.Preload("User").
		Count(&total).
		Scopes(Paginate(page, pageSize)).
		Find(&list)

	return list, total, res.Error
}

/**
 * GetDetail 博客详情
 * @param id {uint} 博客id
 * @return {*Article} 博客
 */
func GetDetail(db *gorm.DB, id uint) (data *Article, err error) {
	res := db.Preload("User").
		Where("parent_id IS NULL AND ID = ?", id).
		First(&data)

	return data, res.Error
}

/**
 * UpdateReplyOrLikeNum 更新博客的回复数或点赞数
 * @param operation {string} addReply:增加回复数 addLike:增加点赞数 subLike:减少点赞数
 */
func UpdateReplyOrLikeNum(db *gorm.DB, id uint, operation string) error {
	var article Article
	db = db.First(&article, id)

	if operation == "addReply" {
		db = db.Update("reply_num", gorm.Expr("reply_num + ?", 1))
	}

	if operation == "addLike" {
		db = db.Update("like_num", gorm.Expr("like_num + ?", 1))
	}

	if operation == "subLike" {
		db = db.Update("like_num", gorm.Expr("like_num - ?", 1))
	}

	return db.Error
}

/**
 * GetReplyList 博客的回复列表(需要parentId)
 * @param id {uint} 博客id
 * @return {[]Article} 回复列表
 * @return {int64} 统计总数
 */
func GetReplyList(db *gorm.DB, page, pageSize int, id uint) (list []Article, total int64, err error) {
	db = db.Model(&Article{})
	res := db.Preload("User").
		Preload("ToAuth").
		Where("parent_id = ?", id).
		Order("created_at desc").
		Count(&total).
		Scopes(Paginate(page, pageSize)).
		Find(&list)

	return list, total, res.Error
}

/**
 * DeleteArticle 删除博客
 * @param id {uint} 博客id
 */
func DeleteArticle(db *gorm.DB, id, userId uint) error {
	var article Article
	db = db.Where("ID = ? AND user_id = ?", id, userId).First(&article)
	if db.Error != nil {
		return db.Error
	}
	res := db.Delete(&article)
	return res.Error
}
