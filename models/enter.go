package models

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key, AUTO_INCREMENT"`
	AuthorID      int64  `json:"author_id,omitempty"`
	Author        *User  `json:"author" gorm:"foreignkey:AuthorID"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty" gorm:"type:tinyint(1);default:0"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty" gorm:"primary_key"`
	UserID     int64  `json:"user_id,omitempty"`
	User       User   `json:"user" gorm:"foreignkey:UserID"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key"`
	Name          string `json:"name,omitempty" gorm:"unique_index"` // 用户名 唯一
	Password      string `json:"password" gorm:"not null, -"`        // 密码
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty" gorm:"type:tinyint(1);default:0"`
}

type Message struct {
	Id         int64  `json:"id,omitempty" gorm:"primary_key,AUTO_INCREMENT"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

// 消息发送事件
type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty" gorm:"primary_key"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

// 消息推送事件
type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
