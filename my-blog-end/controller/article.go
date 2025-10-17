package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"my_blog_back/common"
	"my_blog_back/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Article struct{}

type ReqArticle struct {
	Title   string `json:"title"`
	Content string `json:"content" binding:"required"`
}

type QueryArticleBody struct {
	PageParams
	UserId          uint   `json:"userId"`
	Key             string `json:"key"`
	DescCreatedTime bool   `json:"descCreatedTime"`
}

type ArticleVo struct {
	ID        uint      `json:"ID"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ParentId  uint      `json:"parentId"`
	AuthId    uint      `json:"authId"`
	AuthName  string    `json:"authName"`
	CreatedAt time.Time `json:"createdTime"`
	LikeNum   int64     `json:"likeNum"`
	ReplyNum  int64     `json:"replyNum"`
	LikeFlag  bool      `json:"likeFlag"`
}

func ToArticleVo(article []model.Article, rdb *redis.Client, userId uint) []ArticleVo {
	articleVo := make([]ArticleVo, 0)
	var likeNum int64
	likeFlag := false
	for _, v := range article {
		likeFlag = false
		rdbSetKey := fmt.Sprintf("articlesLikeHash:%v", v.ID)
		likeNum = model.GetRdbArticleLikeNum(rdb, rctx, rdbSetKey)
		if model.CheckRdbHasLikeData(rdb, rctx, rdbSetKey, userId) {
			likeFlag = true
		}
		articleVo = append(articleVo, ArticleVo{
			ID:        v.ID,
			Title:     v.Title,
			Content:   v.Content,
			ParentId:  v.ParentId,
			AuthId:    v.UserID,
			AuthName:  v.User.Name,
			CreatedAt: v.CreatedAt,
			LikeNum:   likeNum,
			ReplyNum:  v.ReplyNum,
			LikeFlag:  likeFlag,
		})
	}
	return articleVo
}

// @Summary 发博客
// @Description 发博客
// @Tags article
// @Param ReqArticle body ReqArticle true "ReqArticle"
// @Success 200 {object} nil
// @Router /article/add [post]
func (*Article) Add(context *gin.Context) {
	DB := common.GetDB()
	req := ReqArticle{}
	err := context.ShouldBindJSON(&req)
	if err != nil {
		slog.Error("/article/add参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	content := req.Content
	title := req.Title

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	newArticle := model.Article{
		Title:   title,
		Content: content,
		UserID:  user.ID,
	}

	err = model.AddArticle(DB, &newArticle)
	if err != nil {
		ReturnServerError(context, nil, "新增博客失败")
		return
	}

	ReturnSuccess(context, nil)
}

// @Summary 博客列表
// @Description 博客列表
// @Tags article
// @Param QueryArticleBody body QueryArticleBody true "QueryArticleBody"
// @Success 200 {object} PageResult[ArticleVo]
// @Router /common/article/query [post]
func (*Article) Query(context *gin.Context) {
	DB := common.GetDB()
	RDB := common.GetRedis()
	query := QueryArticleBody{}
	err := context.ShouldBindJSON(&query)
	if err != nil {
		slog.Error("/common/article/query参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	article, count, err := model.GetArticleOrReplyList(DB, query.Page, query.PageSize, 0, query.Key, query.DescCreatedTime, 0)
	if err != nil {
		ReturnServerError(context, nil, "查询博客失败")
		return
	}

	user, _ := GetCurrentUserInfo(context)

	ReturnSuccess(context, PageResult[ArticleVo]{
		All:     count,
		Content: ToArticleVo(article, RDB, user.ID),
	})
}

// @Summary 某人的博客列表
// @Description 某人的博客列表，不传userId查询当前用户的
// @Tags article
// @Param QueryArticleBody body QueryArticleBody true "QueryArticleBody"
// @Success 200 {object} PageResult[ArticleVo]
// @Router /article/getSomeoneArticle [post]
func (*Article) GetSomeoneArticle(context *gin.Context) {
	DB := common.GetDB()
	RDB := common.GetRedis()
	query := QueryArticleBody{}
	err := context.ShouldBindJSON(&query)
	if err != nil {
		slog.Error("/article/getSomeoneArticle参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	var userId uint
	if query.UserId > 0 {
		userId = query.UserId
	} else {
		userId = user.ID
	}

	article, count, err := model.GetArticleOrReplyList(DB, query.Page, query.PageSize, 0, "", true, userId)
	if err != nil {
		ReturnServerError(context, nil, "查询博客失败")
		return
	}

	ReturnSuccess(context, PageResult[ArticleVo]{
		All:     count,
		Content: ToArticleVo(article, RDB, user.ID),
	})
}

// @Summary 博客详情
// @Description 博客详情
// @Tags article
// @Param ID path int true "ID"
// @Success 200 {object} ArticleVo
// @Router /common/article/getById/:ID [get]
func (*Article) GetById(context *gin.Context) {
	DB := common.GetDB()
	RDB := common.GetRedis()
	_articleId, err := strconv.ParseUint(context.Param("ID"), 10, 0)
	if err != nil {
		ReturnOtherError(context, nil, "博客ID解析错误")
		return
	}
	articleId := uint(_articleId)

	article, err := model.GetDetail(DB, articleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnServerError(context, nil, "博客飞走啦")
			return
		}
		ReturnServerError(context, nil, "查询博客详情失败")
	}

	user, ok := GetCurrentUserInfo(context)
	likeFlag := false
	if ok {
		rdbSetKey := fmt.Sprintf("articlesLikeHash:%v", articleId)
		if model.CheckRdbHasLikeData(RDB, rctx, rdbSetKey, user.ID) {
			likeFlag = true
		}
	}

	articleVo := ArticleVo{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		ParentId:  article.ParentId,
		AuthId:    article.UserID,
		AuthName:  article.User.Name,
		CreatedAt: article.CreatedAt,
		LikeNum:   article.LikeNum,
		ReplyNum:  article.ReplyNum,
		LikeFlag:  likeFlag,
	}
	ReturnSuccess(context, articleVo)
}

// @Summary 删除博客
// @Description 删除博客(当前用户的)
// @Tags article
// @Param ID path int true "ID"
// @Success 200 {object} nil
// @Router /delete/:ID [delete]
func (*Article) Delete(context *gin.Context) {
	DB := common.GetDB()
	_articleId, err := strconv.ParseUint(context.Param("ID"), 10, 0)
	if err != nil {
		ReturnOtherError(context, nil, "博客ID解析错误")
		return
	}
	articleId := uint(_articleId)

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	err = model.DeleteArticle(DB, articleId, user.ID)
	if err != nil {
		ReturnServerError(context, nil, "无法删除哦")
		return
	}

	ReturnSuccess(context, nil)
}
