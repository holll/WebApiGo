package main

import (
	"WebApiGo/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	router := gin.Default()    //获得路由实例

	//添加中间件
	router.Use(Middleware)
	//注册接口
	wxPushGroup := router.Group("/send")
	{
		wxPushGroup.GET("/", tools.WxPushSendHandler)
	}
	router.POST("/simple/server/post", PostHandler)
	router.PUT("/simple/server/put", PutHandler)
	router.DELETE("/simple/server/delete", DeleteHandler)
	//监听端口
	err := http.ListenAndServe(":8005", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Middleware(c *gin.Context) {
	fmt.Println("this is a middleware!")
}

func PostHandler(c *gin.Context) {
	type JsonHolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	holder := JsonHolder{Id: 1, Name: "my name"}
	//若返回json数据，可以直接使用gin封装好的JSON方法
	c.JSON(http.StatusOK, holder)
	return
}
func PutHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("put success!\n"))
	return
}
func DeleteHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("delete success!\n"))
	return
}
