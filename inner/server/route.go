package server

import (
	"github.com/BAN1ce/Tree/inner/api"
	"github.com/BAN1ce/Tree/inner/controller"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Route(g *gin.Engine, api api.API) {
	var (
		v1             = g.Group("/api/v1")
		baseController = controller.BaseController{Api: api}
	)

	v1.GET("topics/match_topics", controller.GetAllMatchTopics(baseController))

	go func() {
		slog.Info("route init success")
		if err := g.Run(); err != nil {
			panic(err)
		}
	}()
}
