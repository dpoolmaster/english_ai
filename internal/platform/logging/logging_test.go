package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewWritesJSONAndFiltersByLevel(t *testing.T) {
	var output bytes.Buffer
	logger := New(&output, slog.LevelInfo)

	logger.Debug("hidden")
	logger.Info("visible", "request_id", "abc")

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("got %d log lines, want 1: %q", len(lines), output.String())
	}

	var record map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("log line is not JSON: %v", err)
	}
	if record["msg"] != "visible" {
		t.Errorf("msg = %v, want visible", record["msg"])
	}
	if record["level"] != "INFO" {
		t.Errorf("level = %v, want INFO", record["level"])
	}
	if record["request_id"] != "abc" {
		t.Errorf("request_id = %v, want abc", record["request_id"])
	}
}
