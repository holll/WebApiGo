package main

import (
	"WebApiGo/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default() //获得路由实例
	exists, err := tools.PathExists("debug")
	if exists && err == nil {
		gin.SetMode(gin.DebugMode)
		fmt.Println("当前为测试环境")
	} else {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("当前为线上环境")
		//添加中间件
		router.Use(Middleware)
	}

	//注册接口
	wxPushGroup := router.Group("/wxpush")
	{
		wxPushGroup.GET("/send", tools.WxPushSendHandler)
		wxPushGroup.GET("/update", tools.WxPushUpdateHandler)
	}
	//监听端口
	err = http.ListenAndServe(":8005", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Middleware(c *gin.Context) {
	fmt.Println("this is a middleware!")
}
