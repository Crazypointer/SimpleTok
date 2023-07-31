package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	Response
	MessageList []models.Message `json:"message_list"`
}

// 业务逻辑
// 1. 用户A给用户B发送消息
//  1.1 生成chatkey
//  1.2 将聊天信息存入redis 做为缓存 如果redis中没有聊天记录，再从数据库中获取
//  1.3 将聊天信息存入数据库

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	user, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在，请重新登录"})
		return
	}

	userIdB, _ := strconv.Atoi(toUserId)
	chatKey := genChatKey(user.ID, int64(userIdB))

	curMessage := models.Message{
		Content:    content,
		CreateTime: time.Now().Unix(),
		FromUserID: user.ID,
		ToUserID:   int64(userIdB),
	}
	//写入数据库
	tx := global.DB.Begin()
	tx.Create(&curMessage)
	tx.Commit()

	//将message对象转换为json字符串
	curMessageStr, err := json.Marshal(curMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(curMessageStr))
	//将聊天记录存入redis
	global.Redis.LPush(chatKey, curMessageStr)
	c.JSON(http.StatusOK, Response{StatusCode: 0})

}

// MessageChat 聊天记录
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	pmt := c.Query("pre_msg_time")
	fmt.Println(pmt)
	user, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在，请重新登录"})
		return
	}

	userIdB, _ := strconv.Atoi(toUserId)
	chatKey := genChatKey(user.ID, int64(userIdB))

	// 从redis中获取聊天记录
	msgListInRedis, err := global.Redis.LRange(chatKey, 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取聊天记录失败"})
		return
	}

	// 如果redis中没有聊天记录
	if len(msgListInRedis) == 0 {
		// 查询数据库
		var msgListInDB []models.Message
		global.DB.Where("from_user_id = ? and to_user_id = ?", user.ID, userIdB).Or("from_user_id = ? and to_user_id = ?", userIdB, user.ID).Order("create_time desc").Find(&msgListInDB)
		// 将聊天记录存入redis
		for _, msg := range msgListInDB {
			//将message对象转换为json字符串
			curMessageStr, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(curMessageStr))
			//将聊天记录存入redis
			global.Redis.LPush(chatKey, curMessageStr)
		}
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: msgListInDB})
		return
	}
	msgList := make([]models.Message, 0)
	// 将redis中的聊天记录转换为message对象
	for _, msg := range msgListInRedis {
		var tempMsg models.Message
		err := json.Unmarshal([]byte(msg), &tempMsg)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取聊天记录失败"})
			return
		}
		if pmt != "" {
			preMsgTime, _ := strconv.ParseInt(pmt, 10, 64)
			if tempMsg.CreateTime <= preMsgTime {
				continue
			}
		}
		msgList = append(msgList, tempMsg)
	}
	//对聊天记录进行排序
	sort.Slice(msgList, func(i, j int) bool {
		return msgList[i].CreateTime < msgList[j].CreateTime
	})

	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: msgList})

}

// genChatKey 生成聊天的key
func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
