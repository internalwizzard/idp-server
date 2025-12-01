package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func RegisterRoutes(g *echo.Group, secured *echo.Group) {
	h := &Handler{}

	// Public endpoints
	g.GET("/monitor/health", h.Health)
	g.GET("/monitor/ready", h.Ready)

	// Protected endpoints
	monitorGroup := secured.Group("/monitor")
	monitorGroup.GET("/secured", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &struct{ Message string }{Message: "You access!"})
	})
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, &struct{ Message string }{Message: "App is healthy"})
}

func (h *Handler) Ready(c echo.Context) error {
	return c.JSON(http.StatusOK, &struct{ Message string }{Message: "App is ready for receive requests"})
}
