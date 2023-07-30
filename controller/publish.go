package controller

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/service"
	"github.com/Crazypointer/simple-tok/utils"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
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
	finalName := fmt.Sprintf("%d_%d_%s", user.ID, now, filename)
	fmt.Println("finalName:", finalName)

	playUrl := ""
	fmt.Println("playUrl:", playUrl)
	// 计算Hash值
	// 读取文件内容
	f, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	byteData, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoHash := utils.Md5(byteData)
	//TODO: 校验Hash值，数据库中存储Hash值，防止重复上传
	// 查询数据库中是否存在该Hash值
	var video models.Video
	err = global.DB.Where("hash_tag = ?", videoHash).First(&video).Error
	if err == nil {
		// 该视频已经存在
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Video already exists",
		})
		return
	}

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
		AuthorID: user.ID,
		PlayUrl:  playUrl,
		HashTag:  videoHash,
	}
	global.DB.Create(&newVideo)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  playUrl + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	token := c.Query("token") // 用于鉴权
	// 鉴权
	_, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}

	userID := c.Query("user_id")
	var videos []models.Video
	global.DB.Where("author_id", userID).Find(&videos)
	var user models.User
	global.DB.Where("id = ?", userID).First(&user)
	var videoList []Video
	for _, video := range videos {
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
			PlayUrl:       video.PlayUrl,
			Title:         video.Title,
			IsFavorite:    false,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
