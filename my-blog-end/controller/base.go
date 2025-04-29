package controller

import (
	"context"
	"my_blog_back/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// redis context
var rctx = context.Background()

// PageParams分页参数
type PageParams struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// PageResult分页响应数据
type PageResult[T any] struct {
	Page    int   `json:"page_num"`
	Size    int   `json:"page_size"`
	All     int64 `json:"all"`
	Content []T   `json:"content"`
}

func ReturnResponse(context *gin.Context, httpStatus, code int, data any, message string) {
	context.JSON(httpStatus, gin.H{"code": code, "data": data, "message": message})
}

func ReturnSuccess(context *gin.Context, data any) {
	ReturnResponse(context, http.StatusOK, 200, data, "success")
}

func ReturnFail(context *gin.Context, data any, message string) {
	ReturnResponse(context, http.StatusBadRequest, 400, data, message)
}

func ReturnAuthFail(context *gin.Context) {
	ReturnResponse(context, http.StatusUnauthorized, 401, nil, "权限不足")
}

func ReturnServerError(context *gin.Context, data any, message string) {
	ReturnResponse(context, http.StatusInternalServerError, 500, data, message)
}

func ReturnOtherError(context *gin.Context, data any, message string) {
	ReturnResponse(context, http.StatusUnprocessableEntity, 422, data, message)
}

func GetCurrentUserInfo(c *gin.Context) (model.User, bool) {
	user, exists := c.Get("user")
	if exists {
		return user.(model.User), exists
	}
	return model.User{}, exists
}
