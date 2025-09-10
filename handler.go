package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type PingRequest struct {
	Ping string
}
type PingResponse struct {
	Pong string
}

func RegisterRoutes(h *server.Hertz) {
	h.GET("/ping/:id", Ping)
	h.POST("/ping", func(ctx context.Context, c *app.RequestContext) {
		var req PingRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]string{
				"error": "Invalid request",
			})
			return
		}
		hlog.SetLevel(hlog.LevelInfo)
		hlog.CtxInfof(ctx, "sdafasdf")

		c.JSON(consts.StatusOK, map[string]string{
			"message": "User created",
			"pong":    req.Ping,
		})
	})

}

func Ping(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	user, err := ping(id)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, user)
}
func ping(id string) (string, error) {
	return id, nil
}
