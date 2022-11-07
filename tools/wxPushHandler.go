package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
)

func WxPushSendHandler(c *gin.Context) {
	// 以下为必须参数
	uid, _ := c.GetQuery("uid")
	content, _ := c.GetQuery("content")
	if len(uid) == 0 || len(content) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "缺少必须参数"})
	}
	// 以下为可选参数
	//url, _ := c.GetQuery("url")
	//cache, _ := c.GetQuery("cache")
	tmpAccessToken, err := os.ReadFile("a.txt")
	if err != nil {
		panic(err)
	}
	accessToken := string(tmpAccessToken)
	params := url.Values{}
	Url, _ := url.Parse("https://qyapi.weixin.qq.com/cgi-bin/message/send")
	params.Set("access_token", accessToken)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	rep, err := http.Get(urlPath)
	if err != nil || rep.StatusCode != 200 {
		c.JSON(http.StatusOK, gin.H{"code": 501, "msg": "发送失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "发送成功"})
		fmt.Println(RepBodyToStr(rep.Body))
	}
}
