package controller

import "github.com/Crazypointer/simple-tok/models"

var DemoUser = models.User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
