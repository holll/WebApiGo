package tools

import (
	"fmt"
	"strings"
)

var mapFunc map[string]func(map[string]string, string, string) string

func InitMapFunc() map[string]func(map[string]string, string, string) string {
	mapFunc[""] = func(_ map[string]string, message string, userId string) string {
		if strings.Index(message, "BV") == 0 {
			return BiliDown(message, userId)
		}
		return message
	}
	mapFunc["image"] = func(cq map[string]string, message string, _ string) string {
		picMd5 := QQUrlToMd5(cq["url"])
		return PicMd5ToUrl(picMd5)
	}
	mapFunc["video"] = func(cq map[string]string, message string, _ string) string {
		return cq["url"]
	}
	mapFunc["final"] = func(cq map[string]string, message string, _ string) string {
		// 有CQ码，但不是已知类型
		return fmt.Sprintf("暂不支持的CQ码类型：%s", cq["CQ"])
	}

	return mapFunc
}
