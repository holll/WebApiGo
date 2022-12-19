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
		rawData = tools.TransHtmlJson(rawData)
		jsonData := string(rawData)
		postType := gjson.Get(jsonData, "post_type").String()
		if postType == "message" {
			message := gjson.Get(jsonData, "message").String()
			cq := tools.ParseCQ(message)
			if cq["CQ"] == "image" {
				picMd5 := tools.QQUrlToMd5(cq["url"])
				context.JSON(http.StatusOK, gin.H{
					"reply": tools.PicMd5ToUrl(picMd5),
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"reply": message,
				})
			}

		}
	})

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		return
	}
}
