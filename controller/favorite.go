package controller

import (
	"fmt"
	"net/http"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
)

// FavoriteAction 为视频点赞
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	//1. 获取视频id 和 点赞类型
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	if _, exist := usersLoginInfo[token]; exist {
		//当前用户 点赞某个视频
		//2. 根据视频id 获取视频信息
		var video models.Video
		global.DB.Where("id = ?", videoId).First(&video)
		fmt.Println("video:", video)
		// 如果是1 则视频点赞总数 FavoriteCount+1
		// 如果是2 则视频点赞总数 FavoriteCount-1
		if actionType == "1" {
			video.FavoriteCount++
		} else if actionType == "2" {
			video.FavoriteCount--
		}
		//3. 更新视频信息
		global.DB.Save(&video)

		//视频作者的总获赞数+1
		var author models.User
		global.DB.Where("id = ?", video.AuthorID).First(&author)
		author.TotalFavorited++
		global.DB.Save(&author)
		//获取用户信息
		var user models.User
		global.DB.Where("id = ?", usersLoginInfo[token].Id).First(&user)
		// 用户喜欢数+1
		user.FavoriteCount++
		//4. 更新用户信息
		global.DB.Save(&user)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
