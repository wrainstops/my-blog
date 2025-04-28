package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model

	User      User    `gorm:"foreignKey:UserID"`
	UserID    uint    `gorm:"not null;comment:用户id" json:"userId"` // 逻辑外键
	Article   Article `gorm:"foreignKey:ArticleId"`
	ArticleId uint    `gorm:"not null;comment:博客id" json:"articleId"` // 逻辑外键
}

/**
 * GetUserLike 用户所有点赞记录
 * @param id {uint} 用户id
 * @return {[]Like} 点赞
 * @return {int64} 统计总数
 */
func GetUserLike(db *gorm.DB, id uint) (list []Like, total int64, err error) {
	res := db.Model(&Like{}).
		Where("user_id = ?", id).
		Count(&total).
		Find(&list)

	return list, total, res.Error
}

/**
 * CheckHasLikeData 校验用户是否已点赞
 * @param userId {uint} 用户id
 * @param articleId {uint} 博客/回复id
 */
func CheckHasLikeData(db *gorm.DB, userId, articleId uint) error {
	var like Like
	res := db.Where("user_id = ? AND article_id = ?", userId, articleId).First(&like)
	return res.Error
}

/**
 * AddLike 点赞
 */
func AddLike(db *gorm.DB, like *Like) error {
	res := db.Create(&like)
	return res.Error
}

/**
 * DeleteLike 取消点赞
 */
func DeleteLike(db *gorm.DB, userId, articleId uint) error {
	var like Like
	res := db.Where("user_id = ? AND article_id = ?", userId, articleId).Delete(&like)
	return res.Error
}
