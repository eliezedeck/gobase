package logging

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RecoverWithZapLogging is a replacement for the default middleware.Recover() which prints only to echo's default
// logger (which is obviously not Zap).
func RecoverWithZapLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				fields := make([]zapcore.Field, 0, 3)

				// Retrieve the ActivityId
				activityId := c.Get("activityId")
				if activityId != nil {
					fields = append(fields, zap.String("ActivityId", activityId.(string)))
				}
				fields = append(fields, zap.Error(err))

				// Get the stack trace
				stack := make([]byte, 4*1024)
				length := runtime.Stack(stack, true)
				fields = append(fields, zap.String("stack", string(stack[:length])))

				L.Error("[PANIC RECOVER]", fields...)

				c.Error(err)
			}
		}()
		return next(c)
	}
}
