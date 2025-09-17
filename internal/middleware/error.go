package middleware

import (
	"net/http"
	res "omailer/pkg/util/response"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	var errCustom *res.MetaError

	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	switch report.Code {
	case http.StatusNotFound:
		errCustom = res.ErrorBuilder(http.StatusNotFound, err, "not found")
	default:
		errCustom = res.ErrorBuilder(http.StatusInternalServerError, err, "internal server error")
	}

	errCustom.SendError(c)
}
