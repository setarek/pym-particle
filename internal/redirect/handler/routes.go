package handler

import "github.com/labstack/echo/v4"

func (h *Handler) Register(v1 *echo.Group) {
	c := v1.Group("/:shorten")
	c.GET("", h.Redirect)
}
