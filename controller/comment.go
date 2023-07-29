package controller

import (
	"net/http"
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
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")

			var comment models.Comment
			comment.User = user
			comment.Content = text
			comment.CreateDate = time.Now().Format("01-02")

			// 将评论信息写入数据库
			global.DB.Create(&comment)

			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0}, Comment: comment})
			return
		}
		if actionType == "2" {
			commentID := c.Query("comment_id")
			var comment models.Comment
			global.DB.Where("id = ?", commentID).First(&comment)

			c.JSON(http.StatusOK, Response{StatusCode: 0})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录系统!"})
		return
	}
	VideoID := c.Query("video_id")

	// 从数据库中读取评论信息
	var comments []models.Comment
	global.DB.Where("video_id = ?", VideoID).Find(&comments)

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
