package gitlab

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"time"
	"strings"
	"github.com/huangxingx/gitlab-webhook/src/dingding"
	"encoding/json"
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
func (n *notifyMsg) GenContent() string {
	content := fmt.Sprintf("event: %s\nuser: %s\nurl: %s\ntime: %s\nrepo: %s\ndescription: %s",
		n.Event, n.UserName, n.Repository.Url, time.Now().String(), n.Repository.Url, n.Repository.Description)
	return strings.Replace(content, "@", "\\@", -1)
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

	// 通知所有人
	atUser := dingding.AtUser{
		IsAtAll: true,
	}
	msg := notifyMsg{
		Event:      "Push Event",
		UserName:   pushEventReq.UserName,
		Repository: pushEventReq.Repository,
	}

	// 生成通知内容
	newMsg := msg.GenContent()
	body := dingding.NewNotifyReq(newMsg, atUser)

	// 发送通知内容
	result, _ := dingding.SendNotifyToDingding(token, body)
	var resultObj interface{}
	json.Unmarshal(result, resultObj)

	c.JSON(http.StatusOK, resultObj)
}
