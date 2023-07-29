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

// UserResponse 封装请求用户信息的响应
type UserResponse struct {
	Response
	User User `json:"user,omitempty"`
}
