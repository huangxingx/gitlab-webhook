package gitlab

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"strings"
	"io/ioutil"
)

const (
	DingdingBotUrl       = "https://oapi.dingtalk.com/robot/send"
	AccessTokenCustomize = ""
)

const (
	ContentTypeJson = "application/json"
)

// gitlab 事件触发的参数
type pushEventReq struct {
	UserName   string     `json:"user_name"`
	Repository repository `json:"repository"`
	Commits    []commits  `json:"commits"`
}

type repository struct {
	Name        string
	Url         string
	Description string
}

type author struct {
	Name  string
	Email string
}

type commits struct {
	Message   string
	Timestamp string
	Url       string
	Author    author
}

type notifyMsg struct {
	Event      string
	UserName   string
	Timestamp  time.Time
	Repository repository
}

// 生成通知消息内容
func genContent(n notifyMsg) string {
	content := fmt.Sprintf("event: %s\nuser: %s\nurl: %s\ntime: %s\nrepo: %s\ndescription: %s",
		n.Event, n.UserName, n.Repository.Url, time.Now().String(), n.Repository.Url, n.Repository.Description)
	return strings.Replace(content, "@", "\\@", -1)
}

type atUser struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type textReq struct {
	Content string `json:"content"`
}

type notifyReq struct {
	Msgtype string  `json:"msgtype"`
	Text    textReq `json:"text"`
	At      atUser  `json:"at"`
}

func (n notifyReq) String() []byte {
	notifyMsg, e := json.Marshal(n)
	if e != nil {
		panic(e.Error())
	}
	return notifyMsg
}

func NewNotifyReq(content string, at atUser) *notifyReq {
	return &notifyReq{
		Msgtype: "text",
		Text: textReq{
			Content: content,
		},
		At: at,
	}
}

func PushHandler(c *gin.Context) {
	pushEventReq := pushEventReq{}
	err := c.BindJSON(&pushEventReq)

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("auth: %s\n", pushEventReq.UserName)

	// 叮叮自定义bot url
	token := c.Param("token")
	if token == "" {
		token = AccessTokenCustomize
	}
	url := DingdingBotUrl + "?access_token=" + token

	// 通知所有人
	atUser := atUser{
		IsAtAll: true,
	}
	msg := notifyMsg{
		Event:      "Push Event",
		UserName:   pushEventReq.UserName,
		Repository: pushEventReq.Repository,
	}

	// 生成通知内容
	newMsg := genContent(msg)
	body := NewNotifyReq(newMsg, atUser).String()
	// 发送通知内容
	resp, _ := http.Post(url, ContentTypeJson, strings.NewReader(string(body)))
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resBody))

	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
