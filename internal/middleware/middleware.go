package middleware

import (
	"fmt"
	"net/http"
	"omailer/pkg/constant"
	"omailer/pkg/util/validator"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo) {
	e.Use(Context)
	e.Use(
		echoMiddleware.Recover(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowCredentials, echo.HeaderContentSecurityPolicy},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		}),
		echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format:           fmt.Sprintf("\n| %s | Host: ${host} | Time: ${time_custom} | Status: ${status} | LatencyHuman: ${latency_human} | UserAgent: ${user_agent} | RemoteIp: ${remote_ip} | Method: ${method} | Uri: ${uri} |\n", constant.APP),
			CustomTimeFormat: "2006/01/02 15:04:05",
			Output:           os.Stdout,
		}),
		echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
			XFrameOptions: "DENY",
			XSSProtection: "1; mode=block",
		}),
	)
	e.HTTPErrorHandler = ErrorHandler
	e.Validator = &validator.CustomValidator{Validator: validator.NewValidator()}
}
