package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	MiddlewareLogger    = "REQUEST_LOGGER"
	MiddlewareRequestID = "REQUEST_ID"
)

func RouteLoggerMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUuid := uuid.New().String()
		xForwardedFor := c.GetHeader("X-Forwarded-For")

		requestLogger := log.With().
			Str("xForwardedFor", xForwardedFor).
			Str("reqId", requestUuid).
			Logger()

		c.Set(MiddlewareLogger, &requestLogger)
		c.Set(MiddlewareRequestID, requestUuid)

		c.Next()
	}
}

func GetRequestLogger(c *gin.Context) *zerolog.Logger {
	return c.MustGet(MiddlewareLogger).(*zerolog.Logger)
}
