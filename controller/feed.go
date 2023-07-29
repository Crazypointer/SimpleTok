package controller

import (
	"net/http"
	"time"

	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/service"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	StatusCode int32          `json:"status_code"`
	StatusMsg  string         `json:"status_msg,omitempty"`
	NextTime   int64          `json:"next_time,omitempty"`
	VideoList  []models.Video `json:"video_list,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//从数据库中获取视频列表
	videoList := service.GetFeedList()

	c.JSON(http.StatusOK, FeedResponse{
		StatusCode: 0,
		VideoList:  videoList,
		NextTime:   time.Now().Unix(),
	})
}
