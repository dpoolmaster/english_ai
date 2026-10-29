package logging

import (
	"io"
	"log/slog"
)

// New returns a JSON logger writing to w. Log records never include configured
// secrets by construction; callers are responsible for not logging PII.
func New(w io.Writer, level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
