package main

import (
	"WebApiGo/tools"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
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
	tools.InitSeed()
	tools.InitDb()
	port := flag.String("port", "8005", "端口")
	rate := flag.Float64("rate", 0, "每秒最大访问次数")
	flag.Parse()
	router := gin.Default() //获得路由实例
	router.GET("/status", isOk)
	exists, err := tools.PathExists("debug")
	if exists && err == nil {
		gin.SetMode(gin.DebugMode)
		fmt.Println("当前为测试环境")
	} else {
		gin.SetMode(gin.ReleaseMode)
		router.Use(Authorize)
		if *rate > 0 {
			router.Use(RateLimitMiddleware(*rate, 1))
		}
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
	fmt.Println("启动端口为：", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", *port), router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func isOk(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "isOk"})
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

func RateLimitMiddleware(rate float64, capacity int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucketWithRate(rate, capacity)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if !bucket.WaitMaxDuration(1, 3*time.Second) {
			c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "rate limit..."})
			c.Abort()
			return
		}
		c.Next()
	}
}
