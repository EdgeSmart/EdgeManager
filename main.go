package main

import (
	"fmt"
	"net"
	"os"

	_ "github.com/EdgeSmart/EdgeManager/dao"
	"github.com/EdgeSmart/EdgeManager/middleware"
	"github.com/EdgeSmart/EdgeManager/mqserver"
	"github.com/EdgeSmart/EdgeManager/service/cluster"
	"github.com/EdgeSmart/EdgeManager/service/user"
	"github.com/EdgeSmart/EdgeManager/token"
	"github.com/gin-gonic/gin"
)

/*
Restful API
*/

func main() {
	signal := make(chan os.Signal)

	go mqserver.Run()
	go httpServer()
	go proxyServer()
	go httpServerTest()

	quitSignal := <-signal
	fmt.Println("Quit", quitSignal)
}

func proxyServer() {
	service := ":8081"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Println(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go mqserver.HandleClient(conn)
	}
}

func httpServer() {
	app := gin.Default()
	tokenConf := token.Config{}
	token.NewInstance("edge", "memery", tokenConf)
	// app.Use(middleware.LoginControl)
	userGroup := app.Group("/user")
	userGroup.POST("/login", user.Login)
	userGroup.POST("/info", user.Info)

	// cluster
	clusterGroup := app.Group("/cluster")
	clusterGroup.POST("/list", cluster.List)

	app.Run()
}

func httpServerTest() {
	app := gin.Default()
	app.Use(middleware.DockerAPI)
	app.Run("0.0.0.0:8082")
}
