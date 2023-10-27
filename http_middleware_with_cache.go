package core

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func HTTPMiddlewareFromCache(key func(IHTTPContext) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(IHTTPContext)
			var item interface{}
			err := cc.Cache().GetJSON(&item, key(cc))
			if err != nil && !errors.Is(err, redis.Nil) {
				cc.NewError(err, Error{
					Status:  http.StatusInternalServerError,
					Code:    "CACHE_ERROR",
					Message: "cache internal error"})
			}

			if item != nil {
				return c.JSON(http.StatusOK, item)
			}

			return next(c)
		}
	}
}
