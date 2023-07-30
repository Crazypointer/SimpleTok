package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction 关注/取消关注
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUser := c.Query("to_user_id")
	toUserID, _ := strconv.ParseInt(toUser, 10, 64)
	actionType := c.Query("action_type")
	// 鉴权
	user, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请先登录"})
		return
	}
	if user.ID == toUserID {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "不能关注或取关自己"})
		return
	}

	if actionType == "1" {
		// 开启事务
		tx := global.DB.Begin()
		// 是否已经关注
		var userFollowRelation models.UserFollowRelation
		err := global.DB.Where("user_id = ? AND follow_user_id = ?", user.ID, toUserID).First(&userFollowRelation).Error
		if err == nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "已关注"})
			return
		}

		// 判断用户是不是互关
		isFollowEachOther := false
		err = global.DB.Where("user_id = ? AND follow_user_id = ?", toUserID, user.ID).First(&userFollowRelation).Error
		if err == nil {
			isFollowEachOther = true // 互关则为好友
			err := tx.Save(&models.UserFollowRelation{
				UserID:       toUserID,
				FollowUserID: user.ID,
				IsFriend:     true,
			}).Error
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
				return
			}
		}

		// 关注关系表添加
		if err := tx.Create(&models.UserFollowRelation{
			UserID:       user.ID,
			FollowUserID: toUserID,
			IsFriend:     isFollowEachOther,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}

		//关注者数据更新
		user.FollowCount++ //关注总数增加
		err = tx.Save(&user).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}

		//被关注用户粉丝数增加
		var followed models.User
		tx.Where("id = ?", toUserID).First(&followed)
		followed.FollowerCount++ //粉丝数增加
		err = tx.Save(&followed).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		tx.Commit()
		// 更新缓存
		usersLoginInfo[token] = user
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return
	}
	// 取关
	if actionType == "2" {
		// 开启事务
		tx := global.DB.Begin()

		//判断两者是不是好友
		var IsFriend models.UserFollowRelation
		if err := tx.Where("user_id = ? AND follow_user_id = ? AND is_friend = true", toUserID, user.ID).First(&IsFriend).Error; err == nil {
			// 互关则为好友
			err := tx.Model(&IsFriend).Update("is_friend", false).Error
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
				return
			}
		}

		// 删除关注关系表
		var userFollowRelation models.UserFollowRelation
		err := tx.Where("user_id = ? AND follow_user_id = ?", user.ID, toUserID).Delete(&userFollowRelation).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}

		//关注者数据更新
		user.FollowCount-- //关注总数减少
		err = tx.Save(&user).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}

		//被关注用户粉丝数增加
		var followed models.User
		tx.Where("id = ?", toUserID).First(&followed)
		followed.FollowerCount-- //粉丝数增加
		err = tx.Save(&followed).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		tx.Commit()
		// 更新缓存
		usersLoginInfo[token] = user
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return
	}

}

// FollowList 返回用户关注列表
func FollowList(c *gin.Context) {
	tk := c.Query("token")
	_, exist := usersLoginInfo[tk]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请先登录"})
		return
	}
	userID := c.Query("user_id")

	var userFollowRelation []models.UserFollowRelation
	if err := global.DB.Where("user_id = ?", userID).Find(&userFollowRelation).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	var userList []User
	for _, v := range userFollowRelation {
		var user models.User
		if err := global.DB.Where("id = ?", v.FollowUserID).First(&user).Error; err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
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
			IsFollow:        true,
		}
		userList = append(userList, userRes)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList 返回用户粉丝列表
func FollowerList(c *gin.Context) {
	tk := c.Query("token")
	_, exist := usersLoginInfo[tk]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请先登录"})
		return
	}
	userID := c.Query("user_id")
	var userFollowers []models.UserFollowRelation
	// follow_user_id 为被关注者id 以此查找其粉丝
	if err := global.DB.Where("follow_user_id = ?", userID).Find(&userFollowers).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	fmt.Println(userFollowers)
	var userList []User
	for _, v := range userFollowers {
		var follower models.User
		if err := global.DB.Where("id = ?", v.UserID).First(&follower).Error; err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		userRes := User{
			ID:              follower.ID,
			Name:            follower.Name,
			Avatar:          follower.Avatar,
			BackgroundImage: follower.BackgroundImage,
			FavoriteCount:   follower.FavoriteCount,
			FollowCount:     follower.FollowCount,
			FollowerCount:   follower.FollowerCount,
			Signature:       follower.Signature,
			TotalFavorited:  follower.TotalFavorited,
			WorkCount:       follower.WorkCount,
			IsFollow:        true,
		}
		userList = append(userList, userRes)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	tk := c.Query("token")
	_, exist := usersLoginInfo[tk]
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请先登录"})
		return
	}
	userID := c.Query("user_id")
	var friends []models.UserFollowRelation
	if err := global.DB.Where("user_id = ? AND is_friend = true", userID).Find(&friends).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	var userList []User
	for _, v := range friends {
		u := models.User{}
		if err := global.DB.Where("id = ?", v.FollowUserID).First(&u).Error; err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		userRes := User{
			ID:              u.ID,
			Name:            u.Name,
			Avatar:          u.Avatar,
			BackgroundImage: u.BackgroundImage,
			FavoriteCount:   u.FavoriteCount,
			FollowCount:     u.FollowCount,
			FollowerCount:   u.FollowerCount,
			Signature:       u.Signature,
			TotalFavorited:  u.TotalFavorited,
			WorkCount:       u.WorkCount,
		}
		userList = append(userList, userRes)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
