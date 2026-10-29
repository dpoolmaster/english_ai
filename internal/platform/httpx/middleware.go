package httpx

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AccessLog writes one structured record per request. It logs method, matched
// route, status and timing only; raw paths, query strings, headers and bodies
// are never logged.
func AccessLog(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()

		status := c.Writer.Status()
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		attrs := []any{
			"method", c.Request.Method,
			"route", route,
			"status", status,
			"duration_ms", time.Since(started).Milliseconds(),
			"request_id", RequestIDFromContext(c.Request.Context()),
			"trace_id", TraceIDFromContext(c.Request.Context()),
		}

		switch {
		case status >= http.StatusInternalServerError:
			log.Error("http request", attrs...)
		case status >= http.StatusBadRequest:
			log.Warn("http request", attrs...)
		default:
			log.Info("http request", attrs...)
		}
	}
}

// Recovery converts panics into a plain 500 response. The panic value is not
// logged because it may contain user supplied data.
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error("panic recovered",
					"request_id", RequestIDFromContext(c.Request.Context()),
					"trace_id", TraceIDFromContext(c.Request.Context()),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// NoRoute returns a small JSON error for unknown paths without leaking details.
func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

// HealthHandler reports process liveness. It deliberately contains no
// configuration or secret material.
func HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func NewRouter(log *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(RequestIDMiddleware(), AccessLog(log), Recovery(log))
	router.NoRoute(NoRoute())
	router.GET("/healthz", HealthHandler())
	return router
}
