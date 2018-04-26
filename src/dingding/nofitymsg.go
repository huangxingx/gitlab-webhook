package dingding

import (
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
)

const (
	BotUrl               = "https://oapi.dingtalk.com/robot/send"
	AccessTokenCustomize = "32b25ba9d8de9a4bcc27a8bb6e1b5bdc6462f2604f5c29723b52686666664aee"
	ContentTypeJson      = "application/json"
)

type AtUser struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type TextReq struct {
	Content string `json:"content"`
}

type NotifyReq struct {
	Msgtype string  `json:"msgtype"`
	Text    TextReq `json:"text"`
	At      AtUser  `json:"at"`
}

func (n NotifyReq) String() string {
	notifyMsg, e := json.Marshal(n)
	if e != nil {
		panic(e.Error())
	}
	return string(notifyMsg)
}

func NewNotifyReq(content string, at AtUser) NotifyReq {
	return NotifyReq{
		Msgtype: "text",
		Text: TextReq{
			Content: content,
		},
		At: at,
	}
}

// 向叮叮接口发送消息
// token 叮叮机器人  access_token
// req NotifyReq 通知请求
func SendNotifyToDingding(token string, req NotifyReq) ([]byte, error) {
	if token == "" {
		token = AccessTokenCustomize
	}

	url := fmt.Sprintf(BotUrl+"?access_token=%s", token)
	resp, _ := http.Post(url, ContentTypeJson, strings.NewReader(req.String()))
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resBody))
	return resBody, nil
}
