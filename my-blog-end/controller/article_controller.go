package controller

import (
	"errors"
	"log/slog"
	"my_blog_back/common"
	"my_blog_back/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Article struct{}

type ReqArticle struct {
	Title   string `json:"title"`
	Content string `json:"content" binding:"required"`
}

type QueryArticleBody struct {
	PageParams
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

func ToArticleVo(article []model.Article, like []model.Like) []ArticleVo {
	articleVo := make([]ArticleVo, 0)
	likeFlag := false
	for _, v := range article {
		likeFlag = false
		for _, v2 := range like {
			if v.ID == v2.ArticleId {
				likeFlag = true
				break
			}
		}
		articleVo = append(articleVo, ArticleVo{
			ID:        v.ID,
			Title:     v.Title,
			Content:   v.Content,
			ParentId:  v.ParentId,
			AuthId:    v.UserID,
			AuthName:  v.User.Name,
			CreatedAt: v.CreatedAt,
			LikeNum:   v.LikeNum,
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
		slog.Error("Login参数格式有误")
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
	query := QueryArticleBody{}
	err := context.ShouldBindJSON(&query)
	if err != nil {
		return
	}

	article, count, err := model.GetArticleOrReplyList(DB, query.Page, query.PageSize, 0, query.Key, query.DescCreatedTime, 0)
	if err != nil {
		ReturnServerError(context, nil, "查询博客失败")
		return
	}

	user, ok := GetCurrentUserInfo(context)
	like := make([]model.Like, 0)
	if ok {
		like, _, err = model.GetUserLike(DB, user.ID)
		if err != nil {
			ReturnServerError(context, nil, "查询用户点赞失败")
			return
		}
	}

	ReturnSuccess(context, PageResult[ArticleVo]{
		All:     count,
		Content: ToArticleVo(article, like),
	})
}

// @Summary 我的博客列表
// @Description 我的博客列表
// @Tags article
// @Param PageParams body PageParams true "PageParams"
// @Success 200 {object} PageResult[ArticleVo]
// @Router /article/getMyArticle [post]
func (*Article) GetMyArticle(context *gin.Context) {
	DB := common.GetDB()
	query := PageParams{}
	err := context.ShouldBindJSON(&query)
	if err != nil {
		return
	}

	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}
	article, count, err := model.GetArticleOrReplyList(DB, query.Page, query.PageSize, 0, "", true, user.ID)
	if err != nil {
		ReturnServerError(context, nil, "查询博客失败")
		return
	}

	like := make([]model.Like, 0)
	if ok {
		like, _, err = model.GetUserLike(DB, user.ID)
		if err != nil {
			ReturnServerError(context, nil, "查询用户点赞失败")
			return
		}
	}

	ReturnSuccess(context, PageResult[ArticleVo]{
		All:     count,
		Content: ToArticleVo(article, like),
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
		err = model.CheckHasLikeData(DB, user.ID, articleId)
		if err == nil {
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
