package controller

import (
	"net/http"
	"time"

	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/service"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	NextTime  int64          `json:"next_time,omitempty"`
	VideoList []models.Video `json:"video_list,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//从数据库中获取视频列表
	videoList := service.GetFeedList()

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
