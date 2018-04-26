package gitlab

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/huangxingx/gitlab-webhook/src/dingding"
	"encoding/json"
)

type userReq struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type objectAttributesReq struct {
	TargetBranch string `json:"target_branch"`
	SourceBranch string `json:"source_branch"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Description  string `json:"description"`
	Title        string `json:"title"`
}

type AuthorReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type lastCommitReq struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Author    author `json:"author"`
}

// 事件触发请求体
type MergeRequestEventReq struct {
	ObjectKind       string              `json:"object_kind"`
	User             userReq             `json:"user"`
	ObjectAttributes objectAttributesReq `json:"object_attributes"`
	LastCommit       lastCommitReq       `json:"last_commit"`
}


// handler
func MergeRequestEventHandler(c *gin.Context) {
	reqParam := MergeRequestEventReq{}
	c.BindJSON(&reqParam)

	body := dingding.NewNotifyReq("",
		dingding.AtUser{
			IsAtAll: false,
		})
	token := c.Param("token")
	res, _ := dingding.SendNotifyToDingding(token, body)
	var resObj interface{}
	json.Unmarshal(res, resObj)
	c.JSON(http.StatusOK, resObj)
}
