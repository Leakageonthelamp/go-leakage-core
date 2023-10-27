package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTPServer(options *HTTPContextOptions) *echo.Echo {
	e := echo.New()
	e.Use(Core(options))

	if options.ContextOptions.ENV.Config().SentryDSN != "" {
		e.Use(sentryecho.New(sentryecho.Options{
			Repanic: true,
		}))
	}

	// Set up debug logging if log level is debug
	//if options.ContextOptions.ENV.Config().LogLevel == logrus.DebugLevel {
	//	e.Debug = true
	//	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//		Format: "method=${method}, uri=${uri}, status=${status}\n",
	//	}))
	//}

	// Apply secure middleware configurations
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            63072000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
		ContentSecurityPolicy: "frame-ancestors 'self'",
	}))

	// Apply core HTTP middleware
	e.Use(HTTPMiddlewareCORS(options))
	e.Use(HTTPMiddlewareRequestID())
	e.Use(HTTPMiddlewareCreateLogger)
	e.Use(HTTPMiddlewareRecoverWithConfig(options.ContextOptions.ENV, middleware.RecoverConfig{
		StackSize: 1 << 20, // 1 KB
	}))

	// Apply rate limiting middleware
	e.Use(HTTPMiddlewareRateLimit(options))

	// Set the custom error handler and not found handler
	e.HTTPErrorHandler = HTTPMiddlewareHandleError(options.ContextOptions.ENV)
	echo.NotFoundHandler = HTTPMiddlewareHandleNotFound

	// Apply additional secure middleware
	e.Use(middleware.Secure())

	// Hide banner
	e.HideBanner = true

	// Print the HTTP service name
	fmt.Println(fmt.Sprintf("HTTP Service: %s", options.ContextOptions.ENV.Config().Service))

	return e
}

func StartHTTPServer(e *echo.Echo, env IENV) {
	if env.Config().ENV == "dev" {
		// Start server in development mode
		e.Logger.Fatal(e.Start(env.Config().Host))
	} else {
		// Start server in production mode

		// Start server in a goroutine to allow for graceful shutdown
		go func() {
			if err := e.Start(env.Config().Host); err != nil && err != http.ErrServerClosed {
				e.Logger.Fatal("shutting down the server")
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		// Create a context with a timeout of 10 seconds for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the server with the given context
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
