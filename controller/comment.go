package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []models.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment models.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	user, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
		return
	}
	actionType := c.Query("action_type")
	vid := c.Query("video_id")
	videoID, _ := strconv.ParseInt(vid, 10, 64)
	tx := global.DB.Begin()
	var video models.Video
	tx.Where("id = ?", videoID).First(&video)
	if actionType == "1" {
		text := c.Query("comment_text")
		var comment models.Comment
		comment.User = user
		comment.Content = text
		comment.CreateDate = time.Now().Local().Format("2006-01-02 15:04")
		comment.VideoID = videoID
		// 将评论信息写入数据库
		tx.Create(&comment)
		//视频评论数+1
		video.CommentCount++
		tx.Save(&video)
		tx.Commit()
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0}, Comment: comment})
		return
	}
	if actionType == "2" {
		commentID := c.Query("comment_id")
		// 删除评论
		tx.Where("id = ?", commentID).Delete(&models.Comment{})
		// 视频评论数-1
		video.CommentCount--
		tx.Save(&video)
		tx.Commit()
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
		return
	}
	videoID := c.Query("video_id")

	// 从数据库中读取评论信息
	var comments []models.Comment
	// 按时间倒序排列
	global.DB.Where("video_id = ?", videoID).Order("create_date desc").Find(&comments)

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
