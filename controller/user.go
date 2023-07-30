package controller

import (
	"errors"
	"net/http"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]models.User{}

type UserLoginResponse struct {
	Response
	UserID int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// 判断数据库中是否存在该用户
	err := global.DB.Where("name = ?", username).First(&models.User{})
	if err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		//为新用户创建随机头像
		avatar := "https://picsum.photos/200"
		// 用户背景图
		background := "https://picsum.photos/seed/picsum/200/300"

		// 随机生成个性签名
		signature := "这个人很懒，什么都没有留下"

		newUser := models.User{
			Name:            username,
			Password:        password,
			Avatar:          avatar,
			BackgroundImage: background,
			Signature:       signature,
		}
		// 将用户信息存入数据库
		global.DB.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserID:   newUser.ID,
			Token:    username + password,
		})
		// TODO: 将用户信息存入内存 其实需要存入redis
		usersLoginInfo[token] = newUser
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// 判断数据库中是否存在该用户
	var user models.User
	err := global.DB.Where("name = ?", username).First(&user).Error

	if err == nil {
		usersLoginInfo[token] = user
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserID:   user.ID,
			Token:    token,
		})
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userID := c.Query("user_id")
	_, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请先登录"})
		return
	}
	var user models.User

	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	// 判断是否关注
	var follow models.UserFollowRelation
	var isFollow bool
	if err := global.DB.Where("user_id = ? AND follow_user_id = ?", usersLoginInfo[token].ID, userID).First(&follow).Error; err == nil {
		isFollow = true
	} else {
		isFollow = false
	}
	userRes := User{
		ID:              user.ID,
		Name:            user.Name,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		IsFollow:        isFollow,
	}
	c.JSON(http.StatusOK, UserResponse{Response: Response{StatusCode: 0}, User: userRes})
}
