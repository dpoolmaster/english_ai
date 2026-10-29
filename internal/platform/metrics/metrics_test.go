package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMuxServesMetrics(t *testing.T) {
	mux := New().Mux()

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	body := response.Body.String()
	if !strings.Contains(body, "go_goroutines") {
		t.Errorf("metrics body does not contain go_goroutines")
	}
	if strings.Contains(body, "password") {
		t.Errorf("metrics body contains secret-like content")
	}
}
