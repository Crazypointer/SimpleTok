package service

import (
	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
)

func GetFeedList() []models.Video {
	var videoList []models.Video
	global.DB.Find(&videoList)
	for i, video := range videoList {
		// 根据视频ID获取视频作者
		var user models.User
		global.DB.Where("id = ?", video.AuthorID).First(&user)
		videoList[i].Author = user
	}
	return videoList
}
