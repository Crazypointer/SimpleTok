package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// User 为前端定制的User结构体
type User struct {
	ID              int64  `json:"id,omitempty"`     //用户id
	Name            string `json:"name,omitempty"`   // 用户名 唯一
	Password        string `json:"-"`                // 密码
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢作品数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int    `json:"total_favorited"`  // 总获赞数量
	WorkCount       int64  `json:"work_count" `
	IsFollow        bool   `json:"is_follow"` //true-已关注，false-未关注 在返回用户列表时使用，通过查询UserFollowRelation表来赋值返回
}

type Video struct {
	ID            int64  `json:"id"`             // 视频唯一标识
	Title         string `json:"title"`          // 视频标题
	Author        User   `json:"author"`         // 视频作者信息
	CommentCount  int64  `json:"comment_count"`  // 视频的评论总数
	FavoriteCount int64  `json:"favorite_count"` // 视频的点赞总数
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	CoverUrl      string `json:"cover_url"`      // 视频封面地址
	PlayUrl       string `json:"play_url"`       // 视频播放地址
}

// UserResponse 封装请求用户信息的响应
type UserResponse struct {
	Response
	User User `json:"user,omitempty"`
}
