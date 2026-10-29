package httpx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	HeaderRequestID   = "X-Request-ID"
	HeaderTraceID     = "X-Trace-ID"
	HeaderTraceparent = "traceparent"

	maxRequestIDLength = 128
)

type contextKey int

const (
	requestIDKey contextKey = iota
	traceIDKey
)

func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey).(string)
	return value
}

func TraceIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(traceIDKey).(string)
	return value
}

// RequestIDMiddleware accepts a caller supplied request id or generates a new
// one, and derives a trace id from a valid W3C traceparent header. The trace id
// falls back to the request id when no valid traceparent is present.
func RequestIDMiddleware() gin.HandlerFunc {
	var fallbackCounter atomic.Uint64

	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(HeaderRequestID))
		if !validIdentifier(requestID) {
			requestID = newIdentifier(&fallbackCounter)
		}

		traceID := traceIDFromTraceparent(c.GetHeader(HeaderTraceparent))
		if traceID == "" {
			traceID = requestID
		}

		ctx := context.WithValue(c.Request.Context(), requestIDKey, requestID)
		ctx = context.WithValue(ctx, traceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("request_id", requestID)
		c.Set("trace_id", traceID)

		c.Header(HeaderRequestID, requestID)
		c.Header(HeaderTraceID, traceID)
		c.Next()
	}
}

func validIdentifier(value string) bool {
	if value == "" || len(value) > maxRequestIDLength {
		return false
	}
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-', r == '_', r == '.':
		default:
			return false
		}
	}
	return true
}

func newIdentifier(fallback *atomic.Uint64) string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err == nil {
		return hex.EncodeToString(raw[:])
	}
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), fallback.Add(1))
}

func traceIDFromTraceparent(header string) string {
	parts := strings.Split(strings.TrimSpace(header), "-")
	if len(parts) != 4 {
		return ""
	}
	if parts[0] != "00" {
		return ""
	}
	if len(parts[1]) != 32 || !isHex(parts[1]) || strings.Trim(parts[1], "0") == "" {
		return ""
	}
	if len(parts[2]) != 16 || !isHex(parts[2]) {
		return ""
	}
	if len(parts[3]) != 2 || !isHex(parts[3]) {
		return ""
	}
	return parts[1]
}

func isHex(value string) bool {
	for _, r := range value {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return false
		}
	}
	return true
}
