package main

import (
	_ "github.com/EdgeSmart/EdgeManager/dao"
	"github.com/EdgeSmart/EdgeManager/middleware"
	"github.com/EdgeSmart/EdgeManager/service/user"
	"github.com/EdgeSmart/EdgeManager/token"
	"github.com/gin-gonic/gin"
)

/*
Restful API
*/

func main() {
	app := gin.Default()
	tokenConf := token.Config{}
	token.NewInstance("edge", "memery", tokenConf)
	app.Use(middleware.LoginControl)
	userGroup := app.Group("/user")
	userGroup.POST("/login", user.Login)
	userGroup.POST("/test", user.Test)
	app.Run()
}
