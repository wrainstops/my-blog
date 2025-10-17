package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"my_blog_back/common"
	"my_blog_back/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct{}

type Contacts struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// @Summary 获取其他用户的信息
// @Description 获取其他用户的信息
// @Tags user
// @Param userId path int true "用户id"
// @Success 200 {object} UserVo
// @Router /user/getOtherUserInfo/:userId [get]
func (*User) GetOtherUserInfo(context *gin.Context) {
	DB := common.GetDB()
	_userId, err := strconv.ParseUint(context.Param("userId"), 10, 0)
	if err != nil {
		ReturnOtherError(context, nil, "用户Id解析错误")
		return
	}
	userId := uint(_userId)

	userInfo, err := model.GetUserById(DB, userId)
	if err != nil {
		ReturnServerError(context, nil, "查询用户信息错误")
		log.Printf("查询用户信息getOtherUserInfo错误: %v", err)
		return
	}

	ReturnSuccess(context, ToUserVo(userInfo))
}

// @Summary 获取其他用户的统计信息
// @Description 获取其他用户的统计信息
// @Tags user
// @Param userId path int true "用户id"
// @Success 200 {object} StatVo
// @Router /user/getOtherUserStats/:userId [get]
func (*User) GetOtherUserStats(context *gin.Context) {
	DB := common.GetDB()
	_userId, err := strconv.ParseUint(context.Param("userId"), 10, 0)
	if err != nil {
		ReturnOtherError(context, nil, "用户Id解析错误")
		return
	}
	userId := uint(_userId)

	articleNum, likeNum, beRepliedNum, beLikedNum, err := model.GetStat(DB, userId)
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

// @Summary 获取最近联系人
// @Description 获取最近聊天的联系人
// @Tags user
// @Success 200 {object} []Contacts
// @Router /user/getRecentContacts [get]
func (*User) GetRecentContacts(context *gin.Context) {
	RDB := common.GetRedis()
	user, ok := GetCurrentUserInfo(context)
	if !ok {
		ReturnOtherError(context, nil, "获取用户信息错误")
		return
	}

	key := fmt.Sprintf("%d:recentContacts", user.ID)
	recentContacts := model.GetRecentContacts(RDB, context, key)
	contactsList := make([]Contacts, 0)
	for _, value := range recentContacts {
		var contact Contacts
		err := json.Unmarshal([]byte(value), &contact)
		if err != nil {
			ReturnOtherError(context, nil, "解析最近联系人错误")
			log.Printf("解析最近联系人错误: %v\njson: %s", err, value)
			return
		}
		contactsList = append(contactsList, contact)
	}

	ReturnSuccess(context, contactsList)
}
