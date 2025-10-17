package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"my_blog_back/common"
	"my_blog_back/model"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Reply struct{}

type ReqReply struct {
	Content  string `json:"content" binding:"required"`
	ParentId uint   `json:"parentId" binding:"required"`
	ToAuthId uint   `json:"toAuthId" binding:"required"`
}

type QueryReplyBody struct {
	PageParams
	ParentId uint `json:"parentId" binding:"required"`
}

type ReplyVo struct {
	ID         uint      `json:"ID"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ParentId   uint      `json:"parentId"`
	AuthId     uint      `json:"authId"`
	AuthName   string    `json:"authName"`
	ToAuthId   uint      `json:"toAuthId"`
	ToAuthName string    `json:"toAuthName"`
	CreatedAt  time.Time `json:"createdTime"`
	LikeNum    int64     `json:"likeNum"`
	ReplyNum   int64     `json:"replyNum"`
	LikeFlag   bool      `json:"likeFlag"`
}

func ToReplyVo(reply []model.Article, rdb *redis.Client, userId uint) []ReplyVo {
	replyVo := make([]ReplyVo, 0)
	var likeNum int64
	likeFlag := false
	for _, v := range reply {
		likeFlag = false
		rdbSetKey := fmt.Sprintf("articlesLikeHash:%v", v.ID)
		likeNum = model.GetRdbArticleLikeNum(rdb, rctx, rdbSetKey)
		if model.CheckRdbHasLikeData(rdb, rctx, rdbSetKey, userId) {
			likeFlag = true
		}
		// for _, v2 := range like {
		// 	if v.ID == v2.ArticleId {
		// 		likeFlag = true
		// 		break
		// 	}
		// }
		replyVo = append(replyVo, ReplyVo{
			ID:         v.ID,
			Title:      v.Title,
			Content:    v.Content,
			ParentId:   v.ParentId,
			AuthId:     v.UserID,
			AuthName:   v.User.Name,
			ToAuthId:   v.ToAuthId,
			ToAuthName: v.ToAuth.Name,
			CreatedAt:  v.CreatedAt,
			LikeNum:    likeNum,
			ReplyNum:   v.ReplyNum,
			LikeFlag:   likeFlag,
		})
	}
	return replyVo
}

// @Summary 回复
// @Description 回复博客/回复
// @Tags reply
// @Param ReqReply body ReqReply true "ReqReply"
// @Success 200 {object} nil
// @Router /reply/add [post]
func (*Reply) Add(context *gin.Context) {
	DB := common.GetDB()
	req := ReqReply{}
	err := context.ShouldBindJSON(&req)
	if err != nil {
		slog.Error("/reply/add参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	content := req.Content
	parentId := req.ParentId
	toAuthId := req.ToAuthId

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	newReply := model.Article{
		Content:  content,
		UserID:   user.ID,
		ParentId: parentId,
		ToAuthId: toAuthId,
	}
	err = model.AddArticle(DB, &newReply)
	if err != nil {
		ReturnServerError(context, nil, "回复失败")
		return
	}
	ReturnSuccess(context, nil)

	// reply_num++
	err = model.UpdateReplyOrLikeNum(DB, parentId, "addReply")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnOtherError(context, nil, "未找到需要回复的博客")
			return
		}
		ReturnServerError(context, nil, "数据库异常")
		return
	}
}

// @Summary 回复列表
// @Description 博客/回复的回复列表
// @Tags reply
// @Param QueryReplyBody body QueryReplyBody true "QueryReplyBody"
// @Success 200 {object} PageResult[ReplyVo]
// @Router /common/reply/query [post]
func (*Reply) Query(context *gin.Context) {
	DB := common.GetDB()
	RDB := common.GetRedis()
	query := QueryReplyBody{}
	err := context.ShouldBindJSON(&query)
	if err != nil {
		slog.Error("/common/reply/queryparentId空")
		ReturnFail(context, nil, "parentId空")
		return
	}

	reply, count, err := model.GetReplyList(DB, query.Page, query.PageSize, query.ParentId)
	if err != nil {
		ReturnServerError(context, nil, "查询回复列表失败")
		return
	}

	user, _ := GetCurrentUserInfo(context)
	// like := make([]model.Like, 0)
	// if ok {
	// 	like, _, err = model.GetUserLike(DB, user.ID)
	// 	if err != nil {
	// 		ReturnServerError(context, nil, "查询用户点赞失败")
	// 		return
	// 	}
	// }

	ReturnSuccess(context, PageResult[ReplyVo]{
		All:     count,
		Content: ToReplyVo(reply, RDB, user.ID),
	})
}

// @Summary 我的所有回复列表
// @Description 我的所有回复列表
// @Tags reply
// @Param PageParams body PageParams true "PageParams"
// @Success 200 {object} PageResult[ReplyVo]
// @Router /reply/getMyReply [post]
func (*Reply) GetMyReply(context *gin.Context) {
	DB := common.GetDB()
	RDB := common.GetRedis()
	query := PageParams{}
	err := context.ShouldBindJSON(&query)
	if err != nil {
		slog.Error("/reply/getMyReply参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}
	reply, count, err := model.GetArticleOrReplyList(DB, query.Page, query.PageSize, 1, "", true, user.ID)
	if err != nil {
		ReturnServerError(context, nil, "查询我的回复列表失败")
		return
	}

	// like := make([]model.Like, 0)
	// if ok {
	// 	like, _, err = model.GetUserLike(DB, user.ID)
	// 	if err != nil {
	// 		ReturnServerError(context, nil, "查询用户点赞失败")
	// 		return
	// 	}
	// }

	ReturnSuccess(context, PageResult[ReplyVo]{
		All:     count,
		Content: ToReplyVo(reply, RDB, user.ID),
	})
}
