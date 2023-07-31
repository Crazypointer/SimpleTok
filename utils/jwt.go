package utils

import (
	"errors"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/sirupsen/logrus"
)

type JwtPayLoad struct {
	UserID          int64  `json:"id" `              //用户id
	Name            string `json:"name"`             // 用户名 唯一
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢作品数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int    `json:"total_favorited"`  // 总获赞数量
	WorkCount       int64  `json:"work_count"`
}

// CustomClaims 自定义的jwt的负载结构体
type CustomClaims struct {
	JwtPayLoad
	jwt.StandardClaims
}

// GenToken 生成token
func GenToken(user JwtPayLoad) (string, error) {
	var MySecret = []byte(global.Config.Jwt.Secret)
	claim := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * time.Duration(global.Config.Jwt.Expire))), // 过期时间
			Issuer:    global.Config.Jwt.Issuer,                                                    // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(MySecret)
}

// ParseToken 解析token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	var MySecret = []byte(global.Config.Jwt.Secret)
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		logrus.Errorf("解析token失败: %s", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("非法token")
}
