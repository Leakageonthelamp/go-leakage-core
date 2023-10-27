package core

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HTTPMiddlewareRecoverWithConfig(env IENV, config middleware.RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(IHTTPContext)
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if err := recover(); err != nil {
					err, ok := err.(error)
					if !ok {
						err = fmt.Errorf("%v", err)
					}

					if ierr, ok := err.(Error); ok {
						_ = cc.NewError(err, ierr)
					} else {
						_ = cc.NewError(err, Error{
							Status:  http.StatusInternalServerError,
							Code:    "INTERNAL_SERVER_ERROR",
							Message: "Internal server error"})
					}

					cc.Error(err)
				}
			}()
			return next(c)
		}
	}
}

func HTTPMiddlewareHandleError(env IENV) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		errMessage := "Internal server error"
		if env.IsDev() {
			errMessage = err.Error()
		}

		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": errMessage,
		})
	}
}
func HTTPMiddlewareHandleNotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"code":    "URL_NOT_FOUND",
		"message": "url not found",
	})
}
