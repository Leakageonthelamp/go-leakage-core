package core

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
)

func uuid(len int) string {
	bytes := make([]byte, len)
	random.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)[:len]
}

func HTTPMiddlewareRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(IHTTPContext)
			rid := uuid(16)
			if c.Request().Header.Get(echo.HeaderXRequestID) != "" {
				rid = c.Request().Header.Get(echo.HeaderXRequestID)
			}

			c.Set(echo.HeaderXRequestID, rid)
			cc.SetData(echo.HeaderXRequestID, rid)
			c.Request().Header.Set(echo.HeaderXRequestID, rid)
			c.Response().Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}
