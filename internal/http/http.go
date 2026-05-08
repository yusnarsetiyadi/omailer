package http

import (
	"fmt"
	"net/http"
	"omailer/internal/app/omailer"
	"omailer/pkg/constant"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "omailer/docs"
)

// @Summary Home
// @Description Welcome page
// @Tags System
// @Produce plain
// @Success 200 {string} string "Welcome message"
// @Router / [get]
func Init(e *echo.Echo) {

	e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Hello there, welcome to app %s version %s.", constant.APP, constant.VERSION)
		return c.String(http.StatusOK, message)
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	omailer.NewHandler().Route(e.Group("/send"))
}
