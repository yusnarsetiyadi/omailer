package http

import (
	"fmt"
	"net/http"
	"omailer/internal/app/omailer"
	"omailer/pkg/constant"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {

	e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Hello there, welcome to app %s version %s.", constant.APP, constant.VERSION)
		return c.String(http.StatusOK, message)
	})

	omailer.NewHandler().Route(e.Group("/send"))
}
