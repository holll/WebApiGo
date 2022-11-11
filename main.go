package main

import (
	"WebApiGo/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	UserTokenSalt  = "default_salt"
	AdminTokenSalt = "admin_salt"
)

var adminRoute = []string{"/wxpush/send"}

func main() {
	router := gin.Default() //获得路由实例
	exists, err := tools.PathExists("debug")
	if exists && err == nil {
		gin.SetMode(gin.DebugMode)
		fmt.Println("当前为测试环境")
	} else {
		gin.SetMode(gin.ReleaseMode)
		router.Use(Authorize)
		fmt.Println("当前为线上环境")
	}

	//注册接口
	wxPushGroup := router.Group("/wxpush")
	{
		wxPushGroup.GET("/send", tools.WxPushSendHandler)
		wxPushGroup.GET("/update", tools.WxPushUpdateHandler)
		wxPushGroup.GET("/clean", tools.WxPushCleanHandler)
	}
	//监听端口
	err = http.ListenAndServe(":8005", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Authorize(c *gin.Context) {
	t := c.Query("t")         // 时间戳
	token := c.Query("token") // 访问令牌
	timeStamp, _ := strconv.ParseInt(t, 10, 64)
	nowTimeStamp := time.Now().Unix()

	saltToken := UserTokenSalt
	if tools.IsInStrSlice(c.FullPath(), adminRoute) {
		saltToken = AdminTokenSalt
	}
	if strings.ToLower(tools.MD5([]byte(t+saltToken))) == strings.ToLower(token) && math.Abs(float64(nowTimeStamp-timeStamp)) < 30 {
		// 验证通过，会继续访问下一个中间件
		c.Next()
	} else {
		// 验证不通过，不再调用后续的函数处理
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
	}
}
