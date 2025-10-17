package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"my_blog_back/common"
	"my_blog_back/model"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Like struct{}

type ReqLike struct {
	ArticleId uint `json:"articleId" binding:"required"`
}

// @Summary 点赞
// @Description 点赞
// @Tags like
// @Param ReqLike body ReqLike true "ReqLike"
// @Success 200 {object} nil
// @Router /like/add [post]
func (*Like) Add(context *gin.Context) {
	// DB := common.GetDB()
	RDB := common.GetRedis()
	addLike := ReqLike{}
	err := context.ShouldBindJSON(&addLike)
	if err != nil {
		slog.Error("/like/add参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	articleId := addLike.ArticleId

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	rdbSetKey := fmt.Sprintf("articlesLikeHash:%v", articleId)
	if model.CheckRdbHasLikeData(RDB, rctx, rdbSetKey, user.ID) {
		ReturnOtherError(context, nil, "本用户已点赞过本博客")
		return
	} else {
		RDB.SAdd(rctx, rdbSetKey, user.ID)
	}

	ReturnSuccess(context, nil)
}

// @Summary 取消点赞
// @Description 取消点赞
// @Tags like
// @Param ReqLike body ReqLike true "ReqLike"
// @Success 200 {object} nil
// @Router /like/cancel [post]
func (*Like) Cancel(context *gin.Context) {
	DB := common.GetDB()
	RDB := common.GetRedis()
	cancelLike := ReqLike{}
	err := context.ShouldBindJSON(&cancelLike)
	if err != nil {
		slog.Error("/like/cancel参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	articleId := cancelLike.ArticleId

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	rdbSetKey := fmt.Sprintf("articlesLikeHash:%v", articleId)
	if !model.CheckRdbHasLikeData(RDB, rctx, rdbSetKey, user.ID) {
		ReturnOtherError(context, nil, "本用户未点赞过本博客哦~.~")
		return
	} else {
		RDB.SRem(rctx, rdbSetKey, user.ID)
		if RDB.SCard(rctx, rdbSetKey).Val() == 0 {
			err = model.DeleteLike(DB, user.ID, articleId)
			if err != nil {
				ReturnServerError(context, nil, "取消点赞失败")
				return
			}
			// like_num--
			err = model.UpdateReplyOrLikeNum(DB, articleId, "subLike")
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					ReturnOtherError(context, nil, "未找到点赞的博客")
					return
				}
				ReturnServerError(context, nil, "数据库异常")
				return
			}
		}
	}

	ReturnSuccess(context, nil)
}
