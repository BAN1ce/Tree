package controller

import (
	"encoding/base64"
	"github.com/BAN1ce/Tree/inner/api/request"
	"github.com/gin-gonic/gin"
)

func GetAllMatchTopics(controller BaseController) func(g *gin.Context) {

	return func(g *gin.Context) {
		var (
			base64Topic, _ = g.GetQuery("topic")
			topic, err     = base64.StdEncoding.DecodeString(base64Topic)
		)
		if err != nil {
			g.JSONP(400, NewResponse(false, 400, "topic decode error", nil))
			return
		}

		result := controller.Api.GetAllMatchTopics(nil, &request.GetAllMatchTopicsRequest{
			Topic: string(topic),
		})
		g.JSONP(200, NewSuccessResponse(result))
	}
}
