package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/utils"
	"github.com/gin-gonic/gin"
)

type FavoriteListResponse struct {
	Response
	VideoList []models.Video `json:"video_list,omitempty"`
}

// FavoriteAction 为视频点赞
func FavoriteAction(c *gin.Context) {

	_claims, _ := c.Get("claims")
	claims := _claims.(*utils.CustomClaims)

	//1. 获取视频id 和 点赞类型: 1 点赞 2 取消点赞
	vid := c.Query("video_id")
	videoID, _ := strconv.ParseInt(vid, 10, 64)
	actionType := c.Query("action_type")

	userID := claims.UserID

	// 查询视频点赞表 判断用户是否已经点赞 默认没有点赞
	isFavorite := false
	var favorite models.UserFavoriteVideo
	err := global.DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&favorite).Error
	if err == nil {
		//如果有查询到记录，说明用户点赞过
		isFavorite = true
	}

	//2. 获取数据
	tx := global.DB.Begin()
	var video models.Video
	tx.Where("id = ?", videoID).First(&video)

	var author models.User //视频作者
	tx.Where("id = ?", video.AuthorID).First(&author)

	var user models.User //当前登录用户
	tx.Where("id = ?", userID).First(&user)
	//根据点赞类型，进行不同的操作
	if actionType == "1" {
		fmt.Println("点赞: ", isFavorite)
		if isFavorite {
			tx.Commit()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "已经点赞过了，无法重复点赞"})
			return
		}
		video.FavoriteCount++
		//视频作者的总获赞数+1
		author.TotalFavorited++
		// 用户喜欢数+1
		user.FavoriteCount++
		//关联表添加记录
		favorite.UserID = userID
		favorite.VideoID = videoID
		tx.Create(&favorite)
	} else if actionType == "2" {
		if !isFavorite {
			tx.Commit()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "还没有点赞过，无法取消点赞"})
			return
		}
		video.FavoriteCount--
		author.TotalFavorited--
		user.FavoriteCount--
		//删除关联表
		tx.Delete(&favorite, "user_id = ? AND video_id = ?", userID, videoID)
	} else {
		tx.Commit()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "参数错误"})
		return
	}

	//3. 更新信息
	if err := tx.Model(&video).Update("favorite_count", video.FavoriteCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if err := tx.Model(&author).Update("total_favorited", author.TotalFavorited).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if err := tx.Model(&user).Update("favorite_count", user.FavoriteCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	// 获取用户id
	userId := c.Query("user_id")
	userID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	//在关联表中查询用户喜欢的视频id
	var favoriteVideos []models.UserFavoriteVideo
	err = global.DB.Where("user_id = ?", userID).Find(&favoriteVideos).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	// 根据视频id查询视频信息
	var videos []models.Video
	for _, favoriteVideo := range favoriteVideos {
		var video models.Video
		err = global.DB.Where("id = ?", favoriteVideo.VideoID).First(&video).Error
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		//查询视频作者信息
		var author models.User
		err = global.DB.Where("id = ?", video.AuthorID).Find(&author).Error
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		video.Author = author
		videos = append(videos, video)
	}
	c.JSON(http.StatusOK, FavoriteListResponse{Response: Response{StatusCode: 0}, VideoList: videos})
}
