package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/service"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	user := usersLoginInfo[token]
	filename := filepath.Base(data.Filename)
	//获取当前时间
	now := time.Duration(time.Now().UnixNano())
	// 生成文件名
	finalName := fmt.Sprintf("%d_%d_%s", user.Id, now, filename)

	playUrl := ""
	//TODO: 校验Hash值，防止重复上传
	//TODO: 数据库中存储Hash值，防止重复上传
	//TODO: 数据表Vieos中增加hash字段，用于存储Hash值

	// 本地存储
	if global.Config.Local.Enable {
		saveFile := filepath.Join("./public/", finalName)
		if err := c.SaveUploadedFile(data, saveFile); err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		playUrl = global.Config.Server.BaseUrl + "/static/" + finalName
	} else {
		// 上传到COS
		playUrl = service.Upload2Cos(data, finalName)
	}
	// 将视频信息存入数据库
	newVideo := models.Video{
		AuthorID: user.Id,
		PlayUrl:  playUrl,
	}
	global.DB.Create(&newVideo)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  playUrl + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var videoList []models.Video
	token := c.Query("token") // 用于鉴权
	fmt.Println(token)        // 此处还缺对用户进行鉴权操作

	userID := c.Query("user_id")
	global.DB.Where("author_id", userID).Find(&videoList)
	for i, video := range videoList {
		// 根据视频ID获取视频作者
		var user models.User
		global.DB.Where("id = ?", video.AuthorID).First(&user)
		fmt.Println(user)
		videoList[i].Author = &user
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
