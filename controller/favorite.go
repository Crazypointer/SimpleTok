package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FavoriteAction 为视频点赞
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		fmt.Printf("%v", usersLoginInfo[token])
		//当前用户 点赞某个视频
		//1. 获取视频id
		//2. 获取用户id

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
