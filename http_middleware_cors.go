package core

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HTTPMiddlewareCORS(options *HTTPContextOptions) echo.MiddlewareFunc {
	allowOrigins := []string{"*"}

	if options.AllowOrigins != nil && len(options.AllowOrigins) > 0 {
		allowOrigins = options.AllowOrigins
	}

	allowHeaders := []string{
		echo.HeaderOrigin,
		echo.HeaderContentType,
		echo.HeaderAccept,
		echo.HeaderAccessControlAllowOrigin,
		echo.HeaderAccessControlAllowMethods,
		echo.HeaderXRequestID,
		echo.HeaderAuthorization,
		echo.HeaderAuthorization,
		"X-Api-Key",
		"X-Device-Id",
	}

	if options.AllowHeaders != nil && len(options.AllowHeaders) > 0 {
		allowHeaders = append(allowHeaders, options.AllowHeaders...)
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: allowHeaders})
}
