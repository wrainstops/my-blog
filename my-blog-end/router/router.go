package router

import (
	"my_blog_back/controller"
	"my_blog_back/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	authApi    controller.Auth
	userApi    controller.User
	articleApi controller.Article
	replyApi   controller.Reply
	likeApi    controller.Like
)

func CollectRoute(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.ValidateTypeMiddleware())
	AuthRoute(r)
	UserRoute(r)
	ArticleRoute(r)
	ReplyRoute(r)
	LikeRoute(r)
	NeedAuthOrNotRoute(r)
}

// AuthRoute 认证路由 - 注册、登录、获取当前用户信息、当前用户统计项
func AuthRoute(r *gin.Engine) {
	auth := r.Group("/auth")

	auth.POST("/register", authApi.Register)
	auth.POST("/login", authApi.Login)

	auth.Use(middleware.AuthMiddleware())

	auth.GET("/info", authApi.Info)
	auth.GET("/getStats", authApi.GetStats)
}

// UserRoute 用户路由 - 获取其他用户信息、其他用户统计项
func UserRoute(r *gin.Engine) {
	user := r.Group("/user")

	user.Use(middleware.AuthMiddleware())

	user.GET("/getOtherUserInfo/:userId", userApi.GetOtherUserInfo)
	user.GET("/getOtherUserStats/:userId", userApi.GetOtherUserStats)
}

// ArticleRoute 博客路由 - 获取我的博客、创建博客、删除博客
func ArticleRoute(r *gin.Engine) {
	article := r.Group("/article")

	article.Use(middleware.AuthMiddleware())

	article.POST("/add", articleApi.Add)
	article.DELETE("/delete/:ID", articleApi.Delete)
	article.POST("/getSomeoneArticle", articleApi.GetSomeoneArticle)
}

// ReplyRoute 评论路由 - 评论
func ReplyRoute(r *gin.Engine) {
	reply := r.Group("/reply")

	reply.Use(middleware.AuthMiddleware())

	reply.POST("/add", replyApi.Add)
	reply.POST("/getMyReply", replyApi.GetMyReply)
}

// LikeRoute 点赞路由 - 点赞、取消点赞
func LikeRoute(r *gin.Engine) {
	like := r.Group("/like")

	like.Use(middleware.AuthMiddleware())

	like.POST("/add", likeApi.Add)
	like.POST("/cancel", likeApi.Cancel)
}

// NeedAuthOrNotRoute - 博客列表、博客详情、评论列表
func NeedAuthOrNotRoute(r *gin.Engine) {
	com := r.Group("/common")

	com.Use(middleware.AuthMiddlewareNoUser())

	com.POST("/article/query", articleApi.Query)
	com.GET("/article/getById/:ID", articleApi.GetById)
	com.POST("/reply/query", replyApi.Query)
}
