package middleware

import (
	"my_blog_back/controller"

	"github.com/gin-gonic/gin"
)

func ValidateTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err := c.Errors.Last(); err != nil {
			if err.Type == gin.ErrorTypeBind {
				controller.ReturnOtherError(nil, nil, err.Error())
				c.Abort()
			}
		}
	}
}
