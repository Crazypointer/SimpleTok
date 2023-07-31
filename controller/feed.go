package controller

import (
	"net/http"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/utils"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	NextTime  int64   `json:"next_time,omitempty"`
	VideoList []Video `json:"video_list,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {

	token := c.Query("token")
	var user *models.User
	if token != "" {
		claims, err := utils.ParseToken(token)
		//token有效时才执行
		if err == nil {
			userID := claims.UserID
			// 根据用户ID获取用户信息
			global.DB.Where("id = ?", userID).Find(&user)
		}
	}

	//从数据库中获取视频列表
	var videos []models.Video
	global.DB.Find(&videos)
	var videoList []Video
	for _, video := range videos {
		// 根据视频ID获取视频作者
		var author models.User
		global.DB.Where("id = ?", video.AuthorID).Find(&author)
		isFavorite := false
		isfollow := false
		if user != nil {
			// 判断当前用户是否已经点赞了视频
			if err := global.DB.Where("user_id = ? AND video_id = ?", user.ID, video.ID).First(&models.UserFavoriteVideo{}).Error; err == nil {
				isFavorite = true
			}
			// 判断当前用户是否关注了视频作者
			if err := global.DB.Where("user_id = ? AND follow_user_id = ?", user.ID, author.ID).First(&models.UserFollowRelation{}).Error; err == nil {
				isfollow = true
			}
			if user.ID == author.ID {
				isfollow = true
			}
		}

		videoList = append(videoList, Video{
			Author: User{
				ID:              author.ID,
				Name:            author.Name,
				Avatar:          author.Avatar,
				BackgroundImage: author.BackgroundImage,
				FavoriteCount:   author.FavoriteCount,
				FollowCount:     author.FollowCount,
				FollowerCount:   author.FollowerCount,
				Signature:       author.Signature,
				TotalFavorited:  author.TotalFavorited,
				WorkCount:       author.WorkCount,
				IsFollow:        isfollow,
			},
			CommentCount:  video.CommentCount,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			ID:            video.ID,
			IsFavorite:    isFavorite,
			PlayUrl:       video.PlayUrl,
			Title:         video.Title,
		})
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
