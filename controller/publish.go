package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"

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
	_claims, _ := c.Get("claims")
	claims := _claims.(*utils.CustomClaims)

	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	filename := filepath.Base(data.Filename)
	//获取当前时间
	now := time.Duration(time.Now().UnixNano())
	// 生成文件名
	finalName := fmt.Sprintf("%d_%d_%s", claims.UserID, now, filename)
	fmt.Println("finalName:", finalName)
	playUrl := ""
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
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//读取视频文件
	var buf bytes.Buffer
	f, err = data.Open()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	io.Copy(&buf, f)
	// 视频封面
	coverStr, err := GetVideoCover(byteData)

	// 本地存储
	if global.Config.Local.Enable {
		saveFile := filepath.Join("./storage/", finalName)
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
		playUrl = service.Upload2Cos(&buf, finalName)
	}

	// 将视频信息存入数据库
	newVideo := models.Video{
		Title:    title,
		AuthorID: claims.UserID,
		PlayUrl:  playUrl,
		HashTag:  videoHash,
		CoverUrl: coverStr,
	}
	tx := global.DB.Begin()
	tx.Create(&newVideo)
	//用户作品数+1
	var user models.User
	tx.Where("id = ?", claims.UserID).First(&user)
	user.WorkCount++
	tx.Save(&user)
	tx.Commit()
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  playUrl + " uploaded successfully",
	})
}

// PublishList 用户发布的视频列表
func PublishList(c *gin.Context) {
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
func GetVideoCover(videoData []byte) (string, error) {
	// 设置图像文件存放路径
	now := time.Now().Unix()
	imagePath := fmt.Sprintf("./storage/%d.jpg", now)
	// 将视频字节数组写入临时文件
	tmpFile, err := os.CreateTemp("", "temp_video*.mp4")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(videoData); err != nil {
		panic(err)
	}

	imgBuf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(tmpFile.Name()).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(imgBuf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(imgBuf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	if global.Config.Local.Enable {
		err = imaging.Save(img, imagePath)
		if err != nil {
			log.Fatal("生成缩略图失败：", err)
			return "", err
		}
		return global.Config.Server.BaseUrl + fmt.Sprintf("/static/%d.jpg", now), nil
	}
	// cos上传
	coverUrl := service.Upload2Cos(imgBuf, fmt.Sprintf("%d.jpg", now))
	return coverUrl, nil

}
