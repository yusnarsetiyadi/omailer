package omailer

import (
	"github.com/labstack/echo/v4"
)

func (h *handler) Route(v *echo.Group) {
	v.POST("", h.OmailerSend)
	v.GET("/just-message", h.OmailerSendJustMessage)
}
