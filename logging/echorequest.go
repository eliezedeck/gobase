package logging

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLoggerForEcho is a middleware and zap to provide an "access log" logging for each request.
func ZapLoggerForEcho(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Use a logger for this context only, may be replaced if there is an error (see below)
			requestLogger := logger

			start := time.Now()

			err := next(c)
			if err != nil {
				requestLogger = requestLogger.With(zap.Error(err))
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zapcore.Field{
				zap.Int("status", res.Status),
				zap.String("ip", c.RealIP()),
				// zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
				zap.String("time", time.Since(start).String()),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id != "" {
				fields = append(fields, zap.String("request_id", id))
			}

			n := res.Status
			switch {
			case n >= 500:
				requestLogger.Error("Server error", fields...)
			case n >= 400:
				requestLogger.Warn("Client error", fields...)
			case n >= 300:
				requestLogger.Info("Redirection", fields...)
			default:
				requestLogger.Info("Success", fields...)
			}

			return err
		}
	}
}
