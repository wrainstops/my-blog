package middleware

import (
	"my_blog_back/common"
	"my_blog_back/controller"
	"my_blog_back/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 获取authorization header
		tokenString := context.GetHeader("Authorization")

		// validate token format
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			controller.ReturnAuthFail(context)
			context.Abort()
			return
		}

		// 因为 “Bearer ”一共占7位，所以从第7位截取
		tokenString = tokenString[7:]

		// 调用ParseToken方法，在tokenString中解析出claims
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { // 如果有error或token无效
			controller.ReturnAuthFail(context)
			context.Abort()
			return
		}

		// 验证通过后获取claims中的userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 验证用户是否存在
		if userId == 0 {
			controller.ReturnAuthFail(context)
			context.Abort()
			return
		}

		// 用户存在，将user信息写入上下文
		context.Set("user", user)
		context.Next()
	}
}

func AuthMiddlewareNoUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			context.Next()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			context.Next()
			return
		}

		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		if userId == 0 {
			context.Next()
			return
		}

		context.Set("user", user)
		context.Next()
	}
}
