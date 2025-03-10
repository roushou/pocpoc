package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func bindHealthRouter(router *echo.Group) {
	group := router.Group("/_health")
	group.GET("", healthCheck)
}

func healthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}
