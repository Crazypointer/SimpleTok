package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
)

// FavoriteAction 为视频点赞
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	//1. 获取视频id 和 点赞类型
	vid := c.Query("video_id")
	videoID, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	actionType := c.Query("action_type")
	if _, exist := usersLoginInfo[token]; exist {

		userID := usersLoginInfo[token].Id
		//2. 获取数据
		var video models.Video
		global.DB.Where("id = ?", videoID).First(&video)

		var author models.User
		global.DB.Where("id = ?", video.AuthorID).First(&author)

		var user models.User
		err := global.DB.Where("id = ?", userID).First(&user).Error
		fmt.Printf("user:%+v\n", user)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		// 判断用户是否已经点赞 默认没有点赞
		isFavorite := false
		//查询关联表
		var favorite models.UserFavoriteVideo
		err = global.DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&favorite).Error
		if err == nil {
			//如果有查询到记录，说明用户点赞过
			isFavorite = true
		}
		//根据点赞类型，进行不同的操作
		if actionType == "1" && !isFavorite {
			video.FavoriteCount++
			//视频作者的总获赞数+1
			author.TotalFavorited++
			// 用户喜欢数+1
			user.FavoriteCount++
			//创建关联表
			favorite.UserID = userID
			favorite.VideoID = videoID
			global.DB.Create(&favorite)
		} else if actionType == "2" && isFavorite {
			video.FavoriteCount--
			author.TotalFavorited--
			user.FavoriteCount--
			//删除关联表
			global.DB.Delete(&favorite, "user_id = ? AND video_id = ?", userID, videoID)
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			return
		}
		//3. 更新信息
		global.DB.Save(&video)
		global.DB.Save(&author)
		global.DB.Save(&user)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//1. 获取token
	token := c.Query("token")
	userId := c.Query("user_id")
	if _, exist := usersLoginInfo[token]; exist {
		//3. 根据用户信息 获取用户喜欢的视频列表
		var videos []models.Video
		global.DB.Where("author_id = ?", userId).Find(&videos)

		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: videos,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
	}
}
