package service

import (
	"fmt"

	"github.com/RaymondCode/simple-tok/global"
	"github.com/RaymondCode/simple-tok/models"
)

func GetFeedList() []models.Video {
	var videoList []models.Video
	global.DB.Find(&videoList)
	for i, video := range videoList {
		// 根据视频ID获取视频作者
		var user models.User
		global.DB.Where("id = ?", video.AuthorID).First(&user)
		fmt.Println(user)
		videoList[i].Author = &user
	}
	return videoList
}
