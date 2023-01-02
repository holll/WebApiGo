package tools

import (
	"fmt"
	"net/http"
	"net/url"
)

var client = http.Client{
	Transport: &http.Transport{DisableKeepAlives: true},
}

func SendMsgPri(userId, msg string) {
	fullUrl := fmt.Sprintf("%s?user_id=%s&message=%s", SendMsgPriApi, userId, url.QueryEscape(msg))
	fmt.Println(fullUrl)
	client.Get(fullUrl)
}
