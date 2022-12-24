package tools

import (
	"fmt"
	"net/http"
)

var client = http.Client{
	Transport: &http.Transport{DisableKeepAlives: true},
}

func SendMsgPri(userId, msg string) {
	fullUrl := fmt.Sprintf("%s?user_id=%s&message=%s", SendMsgPriApi, userId, msg)
	client.Get(fullUrl)
}
