package tools

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var runTask = make([]string, 0)

func InitSeed() {
	rand.Seed(time.Now().Unix())
}

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//IsNotExist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

func BytesHtmlToJson(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	return data
}

func StrHtmlToJson(data string) string {
	data = strings.Replace(data, "&#91;", "[", -1)
	data = strings.Replace(data, "&#93;", "]", -1)
	data = strings.Replace(data, "&#44;", ",", -1)
	data = strings.Replace(data, "&amp;", "&", -1)
	return data
}

// ParseCQ 解析CQ码，返回map
func ParseCQ(cq string) map[string]string {
	cqMap := make(map[string]string)
	if len(cq) <= 2 {
		return cqMap
	}
	cq = cq[1 : len(cq)-1]
	if strings.Index(cq, "CQ") != 0 {
		return cqMap
	}
	cqSlice := strings.Split(cq, ",")
	cqType := strings.Split(cqSlice[0], ":")
	cqMap[cqType[0]] = cqType[1]
	cqSlice = cqSlice[1:]
	for _, tmpStr := range cqSlice {
		tmpSlice := strings.SplitN(tmpStr, "=", 2)
		tmpKey := tmpSlice[0]
		tmpValue := tmpSlice[1]
		cqMap[tmpKey] = tmpValue
	}
	return cqMap
}

func QQUrlToMd5(url string) string {
	reg := regexp.MustCompile("[a-f\\d]{32}|[A-F\\d]{32}")
	rep := reg.FindAllString(url, -1)
	if len(rep) == 0 {
		return ""
	}
	return rep[0]
}

func PicMd5ToUrl(md5 string) string {
	if len(md5) == 0 {
		return "图片地址解析失败"
	}
	rawUrl := fmt.Sprintf("https://gchat.qpic.cn/gchatpic_new/0/-0-%s/0", md5)
	return rawUrl
}

func strInSlice(strSlice []string, str string) bool {
	if len(strSlice) == 0 {
		return false
	}
	for _, strS := range strSlice {
		if strS == str {
			return true
		}
	}
	return false
}

func OnMessage(context *gin.Context) {
	dataReader := context.Request.Body
	rawData, _ := io.ReadAll(dataReader)
	rawData = BytesHtmlToJson(rawData)
	jsonData := string(rawData)
	jsonData = StrHtmlToJson(jsonData)
	postType := gjson.Get(jsonData, "post_type").String()
	userId := gjson.Get(jsonData, "user_id").String()
	if postType == "message" {
		message := gjson.Get(jsonData, "message").String()
		cq := ParseCQ(message)
		switch cq["CQ"] {
		case "image":
			picMd5 := QQUrlToMd5(cq["url"])
			context.JSON(http.StatusOK, gin.H{
				"reply": PicMd5ToUrl(picMd5),
			})
		case "video":
			context.JSON(http.StatusOK, gin.H{
				"reply": cq["url"],
			})
		//case "":
		//	context.JSON(http.StatusOK, gin.H{
		//		"reply": message,
		//	})
		default:
			if strings.Index(message, "BV") == 0 {
				var ban bool
				msg := "已经提交下载任务"
				if strInSlice(runTask, message) {
					ban = true
					msg = "该任务下载中"
				}
				if len(runTask) != 0 {
					ban = true
					msg = fmt.Sprintf("有正在运行的任务：%s", runTask[0])
				}
				context.JSON(http.StatusOK, gin.H{
					"reply": msg,
				})
				if ban {
					return
				}
				runTask = append(runTask, message)
				go func() {
					defer cleanTask()
					BBDownPath := "/opt/BBDown/BBDown"
					if runtime.GOOS == "windows" {
						BBDownPath = "E:\\重装系统\\常用工具\\BBDown.exe"
					}
					_, err := Command("/", BBDownPath, message)
					if err != nil {
						SendMsgPri(userId, fmt.Sprintf("下载失败：%s", err.Error()))
					} else {
						SendMsgPri(userId, fmt.Sprintf("下载成功：%s", message))
					}
				}()
			} else {
				context.JSON(http.StatusOK, gin.H{
					"reply": fmt.Sprintf("暂不支持的CQ码类型：%s", cq["CQ"]),
				})
			}

		}
	} else if postType == "notice" {
		noticeType := gjson.Get(jsonData, "notice_type").String()
		if noticeType == "offline_file" {
			fileUrl := gjson.Get(jsonData, "file.url").String()
			SendMsgPri(userId, fileUrl)
		}
	}
}

func cleanTask() {
	runTask = make([]string, 0)
}

func Command(path, execPath string, arg ...string) (msg string, err error) {
	var cmd *exec.Cmd
	name := "/bin/bash"
	c := "-c"
	// 根据系统设定不同的命令name
	if runtime.GOOS == "windows" {
		name = "cmd"
		c = "/C"
		arg = append([]string{c, execPath}, arg...)
		cmd = exec.Command(name, arg...)
	} else {
		cmd = exec.Command(execPath, arg...)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Dir = path
	err = cmd.Run()
	log.Println(cmd.Args)
	if err != nil {
		msg = fmt.Sprint(err) + ": " + stderr.String()
		err = errors.New(msg)
		//log.Println("err", err.Error(), "cmd", cmd.Args)
	}
	//log.Println(out.String())
	return
}
