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

			fields := make([]zapcore.Field, 0, 10)

			// Log the ActivityId
			activityId := c.Get("activityId")
			if activityId != nil {
				fields = append(fields, zap.String("ActivityId", activityId.(string)))
			}

			fields = append(fields, zap.Int("status", res.Status))
			fields = append(fields, zap.String("ip", c.RealIP()))
			fields = append(fields, zap.String("host", req.Host))
			fields = append(fields, zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)))
			fields = append(fields, zap.Int64("size", res.Size))
			fields = append(fields, zap.String("user_agent", req.UserAgent()))
			fields = append(fields, zap.String("time", time.Since(start).String()))

			id := req.Header.Get(echo.HeaderXRequestID)
			if id != "" {
				fields = append(fields, zap.String("request_id", id))
			}

			n := res.Status
			switch {
			case n >= 500:
				requestLogger.Error("Server error: {request}", fields...)
			case n >= 400:
				requestLogger.Warn("Client error: {request}", fields...)
			case n >= 300:
				requestLogger.Info("Redirection: {request}", fields...)
			default:
				requestLogger.Info("Success: {request}", fields...)
			}

			return err
		}
	}
}
