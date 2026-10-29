package httpx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func testLogger(output *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func testRouter(log *slog.Logger) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware(), AccessLog(log), Recovery(log))
	router.GET("/context", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"request_id": RequestIDFromContext(c.Request.Context()),
			"trace_id":   TraceIDFromContext(c.Request.Context()),
		})
	})
	return router
}

func TestRequestIDGeneratedAndPropagated(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		wantSameID bool
	}{
		{name: "generated", wantSameID: false},
		{name: "propagated", header: "client-request_1", wantSameID: true},
		{name: "invalid replaced", header: "bad id!", wantSameID: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var logs bytes.Buffer
			router := testRouter(testLogger(&logs))

			request := httptest.NewRequest(http.MethodGet, "/context", nil)
			if test.header != "" {
				request.Header.Set(HeaderRequestID, test.header)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", response.Code)
			}

			responseID := response.Header().Get(HeaderRequestID)
			if responseID == "" {
				t.Fatal("response is missing request id")
			}
			if test.wantSameID && responseID != test.header {
				t.Errorf("request id = %q, want %q", responseID, test.header)
			}
			if !test.wantSameID && responseID == test.header {
				t.Errorf("request id = %q, want a replacement", responseID)
			}
			if response.Header().Get(HeaderTraceID) == "" {
				t.Error("response is missing trace id")
			}

			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["request_id"] != responseID {
				t.Errorf("context request id = %q, want %q", body["request_id"], responseID)
			}
		})
	}
}

func TestTraceIDFromTraceparent(t *testing.T) {
	const traceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	const traceparent = "00-" + traceID + "-00f067aa0ba902b7-01"

	var logs bytes.Buffer
	router := testRouter(testLogger(&logs))

	request := httptest.NewRequest(http.MethodGet, "/context", nil)
	request.Header.Set(HeaderTraceparent, traceparent)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["trace_id"] != traceID {
		t.Errorf("trace id = %q, want %q", body["trace_id"], traceID)
	}
}

func TestInvalidTraceparentFallsBackToRequestID(t *testing.T) {
	var logs bytes.Buffer
	router := testRouter(testLogger(&logs))

	request := httptest.NewRequest(http.MethodGet, "/context", nil)
	request.Header.Set(HeaderTraceparent, "garbage")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["trace_id"] != body["request_id"] {
		t.Errorf("trace id = %q, want request id %q", body["trace_id"], body["request_id"])
	}
}

func TestHealthz(t *testing.T) {
	var logs bytes.Buffer
	router := NewRouter(testLogger(&logs))

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if got := response.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want ok", body["status"])
	}
	if strings.Contains(response.Body.String(), "password") {
		t.Errorf("health response contains secret-like content: %q", response.Body.String())
	}
}

func TestMetricsAreNotPublishedOnUserRouter(t *testing.T) {
	var logs bytes.Buffer
	router := NewRouter(testLogger(&logs))

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
}

func TestUnknownRouteReturnsJSONNotFound(t *testing.T) {
	var logs bytes.Buffer
	router := NewRouter(testLogger(&logs))

	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
	if !strings.Contains(response.Body.String(), "not found") {
		t.Errorf("body = %q, want not found", response.Body.String())
	}
}

func TestRecoveryReturns500WithoutLeakingPanicValue(t *testing.T) {
	var logs bytes.Buffer
	logger := testLogger(&logs)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware(), AccessLog(logger), Recovery(logger))
	router.GET("/panic", func(*gin.Context) {
		panic("user supplied secret value")
	})

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", response.Code)
	}
	if strings.Contains(logs.String(), "user supplied secret value") {
		t.Fatalf("logs leak panic value: %q", logs.String())
	}
	if strings.Contains(response.Body.String(), "user supplied secret value") {
		t.Fatalf("response leaks panic value: %q", response.Body.String())
	}
}

func TestAccessLogDoesNotIncludeQueryString(t *testing.T) {
	var logs bytes.Buffer
	router := testRouter(testLogger(&logs))

	request := httptest.NewRequest(http.MethodGet, "/context?token=super-secret", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if strings.Contains(logs.String(), "super-secret") {
		t.Fatalf("access log leaks query string: %q", logs.String())
	}
}
