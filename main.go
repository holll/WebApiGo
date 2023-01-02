package main

import (
	"WebApiGo/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	tools.InitSeed()

	port := "9091"
	router := gin.Default() //获得路由实例
	exists, err := tools.PathExists("debug")
	if exists && err == nil {
		gin.SetMode(gin.DebugMode)
		fmt.Println("当前为测试环境")
	} else {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("当前为线上环境")
	}

	//监听端口
	fmt.Println("启动端口为：", port)

	router.POST("/", tools.OnMessage)

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("端口已被占用")
		return
	}
}
