package common

import (
	"context"
	"log"
	"my_blog_back/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SyncRdbToMysql(DB *gorm.DB, RDB *redis.Client) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	ctx := context.Background()
	article := model.Article{}
	log.Println("开始同步数据")

Loop:
	for {
		select {
		case <-ticker.C:
			val, _, err := RDB.Scan(ctx, 0, "articlesLikeHash:*", 0).Result()
			if err != nil {
				break Loop
			}
			for _, v := range val {
				articleId, err := strconv.Atoi(v[17:]) // 截取出articleId
				if err != nil {
					continue
				}
				count := RDB.SCard(ctx, v).Val()
				if DB.First(&article, uint(articleId)).Error != nil {
					continue
				}
				DB.Debug().First(&article, uint(articleId)).Update("like_num", count) // 更新article表的like_num字段

				list := RDB.SMembers(ctx, v).Val()
				userIdListLikeThisArticle, _, _ := model.GetArticleLike(DB, uint(articleId)) // 所有点赞过本博客的数据
				log.Printf("userIdListLikeThisArticle: %v\n", userIdListLikeThisArticle)

				dbInRdbMap := make(map[int]bool) // db数据是否在rdb中的map

				for _, id := range list {
					userId, _ := strconv.Atoi(id)
					for _, v2 := range userIdListLikeThisArticle {
						intUserId := int(v2.UserID)
						// log.Printf("userId: %v\n", intUserId)
						dbInRdbMap[intUserId] = false
						if intUserId == userId {
							dbInRdbMap[intUserId] = true
						}
					}
					// log.Printf("dbInRdbMap: %+v\n", dbInRdbMap)
					if dbInRdbMap[userId] {
						continue
					}
					log.Printf("userId: %v; articleId: %v\n", userId, articleId)
					newLike := model.Like{
						UserID:    uint(userId),
						ArticleId: uint(articleId),
					}
					model.AddLike(DB, &newLike) // like表插入数据
				}

				for k, v3 := range dbInRdbMap {
					if !v3 {
						// 删除db中数据
						model.DeleteLike(DB, uint(k), uint(articleId))
					}
				}
			}
		}
	}
}
