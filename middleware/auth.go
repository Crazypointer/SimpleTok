package middleware

import (
	"net/http"

	"github.com/Crazypointer/simple-tok/controller"
	"github.com/Crazypointer/simple-tok/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断token是否存在
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			c.JSON(http.StatusOK, controller.Response{StatusCode: 1, StatusMsg: "未携带token,请登录"})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, controller.Response{StatusCode: 1, StatusMsg: "token解析失败,请重新登录"})
			c.Abort()
			return
		}
		//将用户信息写入上下文
		c.Set("claims", claims)
	}
}
