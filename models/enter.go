package models

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key, AUTO_INCREMENT"`        // 视频唯一标识
	Title         string `json:"title"`                                                  // 视频标题
	AuthorID      int64  `json:"author_id,omitempty"`                                    // 视频作者id
	Author        *User  `json:"author" gorm:"foreignkey:AuthorID"`                      // 视频作者信息
	PlayUrl       string `json:"play_url,omitempty"`                                     // 视频播放地址
	CoverUrl      string `json:"cover_url,omitempty"`                                    // 视频封面地址
	FavoriteCount int64  `json:"favorite_count,omitempty"`                               // 视频的点赞总数
	CommentCount  int64  `json:"comment_count,omitempty"`                                // 视频的评论总数
	IsFavorite    bool   `json:"is_favorite,omitempty" gorm:"type:tinyint(1);default:0"` // true-已点赞，false-未点赞
	// 视频标题
}

type Comment struct {
	Id         int64  `json:"id,omitempty" gorm:"primary_key"`
	UserID     int64  `json:"user_id,omitempty"`
	User       User   `json:"user" gorm:"foreignkey:UserID"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id              int64  `json:"id,omitempty" gorm:"primary_key"`                      //用户id
	Name            string `json:"name,omitempty" gorm:"unique_index"`                   // 用户名 唯一
	Password        string `json:"password" gorm:"not null, -"`                          // 密码
	IsFollow        bool   `json:"is_follow,omitempty" gorm:"type:tinyint(1);default:0"` //true-已关注，false-未关注
	Avatar          string `json:"avatar"`                                               // 用户头像
	BackgroundImage string `json:"background_image"`                                     // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`                                       // 喜欢数
	FollowCount     int64  `json:"follow_count"`                                         // 关注总数
	FollowerCount   int64  `json:"follower_count"`                                       // 粉丝总数
	Signature       string `json:"signature"`                                            // 个人简介
	TotalFavorited  string `json:"total_favorited"`                                      // 获赞数量
	WorkCount       int64  `json:"work_count"`                                           // 作品数
}

type Message struct {
	Id         int64  `json:"id,omitempty" gorm:"primary_key,AUTO_INCREMENT"` // 消息id
	Content    string `json:"content,omitempty"`                              // 消息内容
	CreateTime string `json:"create_time,omitempty"`                          // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID int64  `json:"from_user_id"`                                   // 消息发送者id
	ToUserID   int64  `json:"to_user_id"`                                     // 消息接收者id
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
