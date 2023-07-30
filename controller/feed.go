package controller

import (
	"net/http"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	NextTime  int64   `json:"next_time,omitempty"`
	VideoList []Video `json:"video_list,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//从数据库中获取视频列表
	var videos []models.Video
	global.DB.Find(&videos)
	var videoList []Video
	for _, video := range videos {
		// 根据视频ID获取视频作者
		var user models.User
		global.DB.Where("id = ?", video.AuthorID).First(&user)
		videoList = append(videoList, Video{
			Author: User{
				ID:              user.ID,
				Name:            user.Name,
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				FavoriteCount:   user.FavoriteCount,
				FollowCount:     user.FollowCount,
				FollowerCount:   user.FollowerCount,
				Signature:       user.Signature,
				TotalFavorited:  user.TotalFavorited,
				WorkCount:       user.WorkCount,
				IsFollow:        false,
			},
			CommentCount:  video.CommentCount,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			ID:            video.ID,
			IsFavorite:    false, // 未登陆的用户，视频列表中的视频都是未点赞的
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
