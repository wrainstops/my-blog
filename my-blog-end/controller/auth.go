package controller

import (
	"errors"
	"log"
	"log/slog"
	"my_blog_back/common"
	"my_blog_back/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct{}

type ReqUser struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserVo struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}

type StatVo struct {
	// 发布博客数量、点赞量、被回复量、被点赞量
	ArticleNum   int64 `json:"articleNum"`
	LikeNum      int64 `json:"likeNum"`
	BeRepliedNum int64 `json:"beRepliedNum"`
	BeLikedNum   int64 `json:"beLikedNum"`
}

func ToUserVo(user model.User) UserVo {
	return UserVo{
		ID:   user.ID,
		Name: user.Name,
	}
}

// @Summary 注册
// @Description 用户注册
// @Tags auth
// @Param ReqUser body ReqUser true "ReqUser"
// @Success 200 {object} nil
// @Router /auth/register [post]
func (*Auth) Register(context *gin.Context) {
	DB := common.GetDB()
	req := ReqUser{}
	err := context.ShouldBindJSON(&req)
	if err != nil {
		slog.Error("Register参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	// 获取参数
	name := req.Name
	password := req.Password

	// 数据验证
	if len(password) < 6 {
		// 422状态码，接受表单数据，但是数据校验不过，无法处理提交过来的数据
		ReturnOtherError(context, nil, "密码长度必须大于等于6")
		return
	}

	//新建用户
	//密码加密存储
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ReturnServerError(context, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:     name,
		Password: string(hashPassword),
		Articles: nil,
	}

	err = model.AddUser(DB, newUser)
	if err != nil {
		ReturnServerError(context, nil, "注册失败")
		return
	}

	ReturnSuccess(context, nil)
}

// @Summary 登录
// @Description 用户登录
// @Tags auth
// @Param ReqUser body ReqUser true "ReqUser"
// @Success 200 {string} token
// @Router /auth/login [post]
func (*Auth) Login(context *gin.Context) {
	DB := common.GetDB()
	req := ReqUser{}
	err := context.ShouldBindJSON(&req)
	if err != nil {
		slog.Error("Login参数格式有误")
		ReturnFail(context, nil, "参数格式有误")
		return
	}

	// 获取参数
	name := req.Name
	password := req.Password

	// 数据验证
	if len(password) < 6 {
		ReturnOtherError(context, nil, "密码长度必须大于等于6")
		return
	}

	// 判断密码是否正确
	user, err := model.GetUserByName(DB, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnOtherError(context, nil, "未找到用户")
			return
		}
		ReturnServerError(context, nil, "数据库异常")
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ReturnFail(context, nil, "密码错误")
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ReturnServerError(context, nil, "系统异常")
		log.Printf("token generate error: %v !", err)
		return
	}

	// 返回结果
	ReturnSuccess(context, token)
}

// @Summary 用户信息
// @Description 当前用户信息
// @Tags auth
// @Success 200 {object} UserVo
// @Router /auth/info [get]
func (*Auth) Info(context *gin.Context) {
	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	ReturnSuccess(context, ToUserVo(user))
}

// @Summary 统计
// @Description 用户的统计信息
// @Tags auth
// @Success 200 {object} StatVo
// @Router /auth/getStats [get]
func (*Auth) GetStats(context *gin.Context) {
	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	DB := common.GetDB()

	articleNum, likeNum, beRepliedNum, beLikedNum, err := model.GetStat(DB, user.ID)
	if err != nil {
		ReturnServerError(context, nil, "查询统计失败")
		return
	}

	result := StatVo{
		ArticleNum:   articleNum,
		LikeNum:      likeNum,
		BeRepliedNum: beRepliedNum,
		BeLikedNum:   beLikedNum,
	}
	ReturnSuccess(context, result)
}
