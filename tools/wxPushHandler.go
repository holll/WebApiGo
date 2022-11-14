package tools

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/url"
	"os"
)

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    wxSend, err := UnmarshalWxSend(bytes)
//    bytes, err = wxSend.Marshal()

func UnmarshalWxSend(data []byte) (WxSend, error) {
	var r WxSend
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *WxSend) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type WxSend struct {
	Touser                 string `json:"touser"`
	Msgtype                string `json:"msgtype"`
	Agentid                string `json:"agentid"`
	Text                   Text   `json:"text"`
	EnableDuplicateCheck   uint8  `json:"enable_duplicate_check"`
	DuplicateCheckInterval uint16 `json:"duplicate_check_interval"`
}

type Text struct {
	Content string `json:"content"`
}

func UnmarshalWxRep(data []byte) (WxRep, error) {
	var r WxRep
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *WxRep) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type WxRep struct {
	Errcode     int64  `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func NewWxSendJson(uid, content, agentId, rawUrl string, check uint8) *WxSend {
	data := new(WxSend)
	data.Touser = uid
	data.Msgtype = "text"
	data.Agentid = agentId
	if len(rawUrl) != 0 {
		content = fmt.Sprintf("%s\n<a href=\"%s\">点击查看原文链接</a>", content, rawUrl)
	}
	data.Text = Text{Content: content}
	data.EnableDuplicateCheck = check
	data.DuplicateCheckInterval = 900
	return data
}

func WxPushSendHandler(c *gin.Context) {
	client := http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}
	// 以下为必须参数
	uid := c.Query("uid")
	content := c.Query("content")
	agentId := c.Query("agentid")
	if len(uid) == 0 || len(content) == 0 || len(agentId) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "缺少必须参数"})
		return
	}
	// 以下为可选参数
	tmpCache := c.Query("cache")
	check := uint8(1)
	if len(tmpCache) != 0 && tmpCache == "0" {
		check = 0
	}
	rawUrl := c.Query("url")
	accessToken, err := SearchToken()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "读取AccessToken失败"})
		return
	}
	accessToken := string(tmpAccessToken)
	data := NewWxSendJson(uid, content, agentId, rawUrl, check)
	marshal, err := data.Marshal()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "参数格式化失败，请检测数据合法性"})
		return
	}
	rep, err := client.Post("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+accessToken, "application/json", bytes.NewReader(marshal))
	if err != nil || rep.StatusCode != 200 {
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 501, "msg": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 501, "msg": "状态码为" + rep.Status})
		}
	} else {
		defer rep.Body.Close()
		wxRep, err := UnmarshalWxRep(RepBodyToByteSlice(rep.Body))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 504, "msg": "请求已发送，但解析响应数据失败"})
			return
		}
		if wxRep.Errcode != 0 {
			c.JSON(http.StatusOK, gin.H{"code": 505, "msg": wxRep.Errmsg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "发送成功"})
	}
}

// WxPushUpdateHandler 必须参数为corpid,corpsecret,agentid
func WxPushUpdateHandler(c *gin.Context) {
	params := url.Values{}
	Url, err := url.Parse("https://qyapi.weixin.qq.com/cgi-bin/gettoken")
	if err != nil {
		return
	}
	corpid := c.Query("corpid")
	corpsecret := c.Query("corpsecret")
	agentid := c.Query("agentid")
	if len(corpid) == 0 || len(corpsecret) == 0 || len(agentid) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "缺少必须参数"})
		return
	}
	params.Set("corpid", corpid)
	params.Set("corpsecret", corpsecret)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	rep, err := http.Get(urlPath)
	if err != nil || rep.StatusCode != 200 {
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 501, "msg": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 501, "msg": "状态码为" + rep.Status})
		}
	} else {
		defer rep.Body.Close()
		wxRep, err := UnmarshalWxRep(RepBodyToByteSlice(rep.Body))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 504, "msg": "请求已发送，但解析响应数据失败"})
			return
		}
		if wxRep.Errcode != 0 {
			c.JSON(http.StatusOK, gin.H{"code": 505, "msg": wxRep.Errmsg})
			return
		}
		accessToken := wxRep.AccessToken
		db, err := sql.Open("sqlite3", "data.db")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
			return
		}
		defer db.Close()

		if !IsNewAgent(corpid, agentid) {
			fmt.Println(fmt.Sprintf("企业：%s，应用：%s，token已存在", corpid, agentid))
			err = UpdateDb(corpsecret, accessToken)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
			}
		} else {
			err = InsertDb(corpid, corpsecret, agentid, accessToken)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
	}
}

func WxPushCleanHandler(c *gin.Context) {
	// Todo 清理数据库，日志清除和accessToken清除
}
