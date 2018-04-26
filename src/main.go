package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huangxingx/gitlab-webhook/src/gitlab"
)

var DB = make(map[string]string)

func setupRouter() *gin.Engine {

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/api/push_event", gitlab.PushHandler)
	r.POST("/api/merge_request_event", gitlab.MergeRequestEventHandler)

	return r

}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
