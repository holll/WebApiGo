package tools

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

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
