package models

type Comment struct {
	ID         int64  `json:"id,omitempty" gorm:"primary_key"`
	UserID     int64  `json:"user_id,omitempty"`
	User       User   `json:"user" gorm:"foreignkey:UserID"`
	VideoID    int64  `json:"video_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type Video struct {
	ID            int64  `json:"id,omitempty" gorm:"primary_key, AUTO_INCREMENT"` // 视频唯一标识
	Title         string `json:"title"`                                           // 视频标题
	HashTag       string `json:"hash_tag,omitempty"`                              // 视频hash tag 避免重复上传
	AuthorID      int64  `json:"author_id,omitempty"`                             // 视频作者id
	Author        User   `json:"author" gorm:"foreignkey:AuthorID"`               // 视频作者信息
	PlayUrl       string `json:"play_url,omitempty"`                              // 视频播放地址
	CoverUrl      string `json:"cover_url,omitempty"`                             // 视频封面地址
	FavoriteCount int64  `json:"favorite_count,omitempty"`                        // 视频的点赞总数
	CommentCount  int64  `json:"comment_count,omitempty"`                         // 视频的评论总数
}

type User struct {
	ID              int64  `json:"id,omitempty" gorm:"primary_key"`    //用户id
	Name            string `json:"name,omitempty" gorm:"unique_index"` // 用户名 唯一
	Password        string `json:"-" gorm:"not null"`                  // 密码
	Avatar          string `json:"avatar"`                             // 用户头像
	BackgroundImage string `json:"background_image"`                   // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count" gorm:"default:0"`    // 喜欢作品数
	FollowCount     int64  `json:"follow_count" gorm:"default:0"`      // 关注总数
	FollowerCount   int64  `json:"follower_count" gorm:"default:0"`    // 粉丝总数
	Signature       string `json:"signature"`                          // 个人简介
	TotalFavorited  int    `json:"total_favorited" gorm:"default:0"`   // 总获赞数量
	WorkCount       int64  `json:"work_count" gorm:"default:0"`
}

// UserFollowRelation 记录用户与用户的对应关系 正向关注 反向粉丝
type UserFollowRelation struct {
	UserID       int64 `gorm:"primary_key"` // 用户id
	FollowUserID int64 // 关注的用户id
	IsFriend     bool  // 是否是好友  互相关注的人互相是好友
}

// UserFavoriteVideo  记录用户id与视频id的对应关系 用来判断是否点赞
type UserFavoriteVideo struct {
	UserID  int64 `gorm:"primary_key"`
	VideoID int64
}

type Message struct {
	ID         int64  `json:"id" gorm:"primary_key,AUTO_INCREMENT"` // 消息id
	CreateTime string `json:"create_time"`                          // 消息发送时间 yyyy-MM-dd HH:MM:ss
	ToUserID   int64  `json:"to_user_id"`                           // 消息接收者id
	FromUserID int64  `json:"from_user_id"`                         // 消息发送者id
	Content    string `json:"content,omitempty"`                    // 消息内容
}

// 消息发送事件表
type MessageSendEvent struct {
	UserID     int64  `json:"user_id,omitempty" gorm:"primary_key"`
	ToUserID   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

// 消息推送事件表
type MessagePushEvent struct {
	FromUserID int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
