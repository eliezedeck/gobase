package logging

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// EchoInjectActivityId adds X-Seq-Activity-Id HTTP header and 'activityId' in the current Echo
// context, useful for logging and tracing.
func EchoInjectActivityId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Add an ActivityId from here and pass it along
		guid, _ := uuid.NewRandom()
		activityId := guid.String()
		c.Set("activityId", activityId)

		// Return back the ActivityId of this request
		c.Response().Header().Set("X-Seq-Activity-Id", activityId)

		return next(c)
	}
}
