package main

import (
	"github.com/Crazypointer/simple-tok/controller"
	"github.com/Crazypointer/simple-tok/middleware"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// 配置cors
	r.Use(middleware.CORSMiddleware())

	// storage directory is used to serve static resources
	r.Static("/static", "./storage")

	// api router group
	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.AuthMiddleware(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middleware.AuthMiddleware(), controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.AuthMiddleware(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.AuthMiddleware(), controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.AuthMiddleware(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", middleware.AuthMiddleware(), controller.FriendList)
	apiRouter.GET("/message/chat/", middleware.AuthMiddleware(), controller.MessageChat)
	apiRouter.POST("/message/action/", middleware.AuthMiddleware(), controller.MessageAction)
}
