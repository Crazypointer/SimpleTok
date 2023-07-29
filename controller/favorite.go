package controller

import (
	"net/http"
	"strconv"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
)

type FavoriteListResponse struct {
	StatusCode int32          `json:"status_code"`
	StatusMsg  string         `json:"status_msg,omitempty"`
	VideoList  []models.Video `json:"video_list,omitempty"`
}

// FavoriteAction 为视频点赞
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
		return
	}
	//1. 获取视频id 和 点赞类型
	vid := c.Query("video_id")
	videoID, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	actionType := c.Query("action_type")

	userID := usersLoginInfo[token].Id
	//2. 获取数据
	var video models.Video
	global.DB.Where("id = ?", videoID).First(&video)

	var author models.User //视频作者
	err = global.DB.Where("id = ?", video.AuthorID).First(&author).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	var user models.User //当前登录用户
	err = global.DB.Where("id = ?", userID).First(&user).Error
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
	//开启事务
	tx := global.DB.Begin()

	// 在事务中执行操作
	if err := tx.Save(&video).Error; err != nil {
		// 发生错误时回滚事务
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})

	} else if err := tx.Save(&author).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		// 提交事务
		tx.Commit()
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//1. 获取token
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
		return
	}
	//2. 获取用户id
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
	//3. 根据视频id查询视频信息
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
		err = global.DB.Where("id = ?", video.AuthorID).First(&author).Error
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		video.Author = author
		videos = append(videos, video)
	}
	c.JSON(http.StatusOK, FavoriteListResponse{StatusCode: 0, VideoList: videos})
}
