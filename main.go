package main

import (
	"WebApiGo/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func main() {
	client := http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}
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

	router.POST("/", func(context *gin.Context) {
		dataReader := context.Request.Body
		rawData, _ := io.ReadAll(dataReader)
		rawData = tools.BytesHtmlToJson(rawData)
		jsonData := string(rawData)
		jsonData = tools.StrHtmlToJson(jsonData)
		postType := gjson.Get(jsonData, "post_type").String()
		if postType == "message" {
			message := gjson.Get(jsonData, "message").String()
			cq := tools.ParseCQ(message)
			if cq["CQ"] == "image" {
				picMd5 := tools.QQUrlToMd5(cq["url"])
				context.JSON(http.StatusOK, gin.H{
					"reply": tools.PicMd5ToUrl(picMd5),
				})
			} else if cq["CQ"] == "video" {
				context.JSON(http.StatusOK, gin.H{
					"reply": cq["url"],
				})
			} else if cq["CQ"] != "" {
				context.JSON(http.StatusOK, gin.H{
					"reply": fmt.Sprintf("暂不支持的CQ码类型：%s", cq["CQ"]),
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"reply": message,
				})
			}
		} else if postType == "notice" {
			noticeType := gjson.Get(jsonData, "notice_type").String()
			userId := gjson.Get(jsonData, "user_id").String()
			if noticeType == "offline_file" {
				fileUrl := gjson.Get(jsonData, "file.url").String()
				// Todo 进一步封装接口
				client.Get(fmt.Sprintf("%s?user_id=%s&message=%s", tools.SendMsgPri, userId, fileUrl))
			}
		}
	})

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		return
	}
}
