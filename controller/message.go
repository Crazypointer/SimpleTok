package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
)

// 临时聊天记录 未来可以存储到redis中
var tempChat = map[string][]models.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []models.Message `json:"message_list"`
}

// MessageAction 聊天记录
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.ID, int64(userIdB))

		atomic.AddInt64(&messageIdSequence, 1)
		curMessage := models.Message{
			ID:         messageIdSequence,
			Content:    content,
			CreateTime: time.Now().Format(time.Kitchen),
		}

		if messages, exist := tempChat[chatKey]; exist {
			tempChat[chatKey] = append(messages, curMessage)
		} else {
			tempChat[chatKey] = []models.Message{curMessage}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all 发送消息
// TODO: 业务逻辑没理清 待完善
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.ID, int64(userIdB))

		var msg models.Message
		text := c.Query("content")
		msg.Content = text
		msg.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		msg.FromUserID = user.ID
		msg.ToUserID = int64(userIdB)

		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// genChatKey 生成聊天的key
func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
