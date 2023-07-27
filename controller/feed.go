package controller

import (
	"net/http"
	"time"

	"github.com/RaymondCode/simple-tok/models"
	"github.com/RaymondCode/simple-tok/service"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response  models.Response `json:"response"`
	VideoList []models.Video  `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//从数据库中获取视频列表
	videoList := service.GetFeedList()

	c.JSON(http.StatusOK, FeedResponse{
		Response:  models.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
