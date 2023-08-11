package controller

import (
	"fmt"
	"net/http"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/utils"
	"github.com/NebulousLabs/fastrand"
	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserID int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 判断数据库中是否存在该用户
	count := global.DB.Where("name = ?", username).First(&models.User{}).RowsAffected
	if count != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		//为新用户创建随机头像
		avatar := GenAvatar()
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
		//生成token
		token, err := utils.GenToken(utils.JwtPayLoad{
			UserID:          newUser.ID,
			Name:            newUser.Name,
			Avatar:          newUser.Avatar,
			BackgroundImage: newUser.BackgroundImage,
			FavoriteCount:   newUser.FavoriteCount,
			FollowCount:     newUser.FollowCount,
			FollowerCount:   newUser.FollowerCount,
			Signature:       newUser.Signature,
			TotalFavorited:  newUser.TotalFavorited,
			WorkCount:       newUser.WorkCount,
		})
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserID:   newUser.ID,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 判断数据库中是否存在该用户
	var user models.User
	err := global.DB.Where("name = ?", username).First(&user).Error

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或密码错误"},
		})
	}

	if user.Password != password {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或密码错误"},
		})
		return
	}
	//生成token
	token, err := utils.GenToken(utils.JwtPayLoad{
		UserID:          user.ID,
		Name:            user.Name,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	})

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserID:   user.ID,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*utils.CustomClaims)

	userID := c.Query("user_id")
	var user models.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	// 判断是否关注
	var follow models.UserFollowRelation
	var isFollow bool
	if err := global.DB.Where("user_id = ? AND follow_user_id = ?", claims.ID, userID).Find(&follow).Error; err == nil {
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
func GenAvatar() string {
	return fmt.Sprintf("https://q.qlogo.cn/headimg_dl?dst_uin=%d&spec=64&img_type=jpg", 200000000+fastrand.Intn(999999999))
}
