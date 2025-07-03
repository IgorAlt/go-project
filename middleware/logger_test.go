package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware_Basic(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest("GET", "/path", nil)
	rec := httptest.NewRecorder()

	LoggerMiddleware(next).ServeHTTP(rec, req)

	if !called {
		t.Errorf("expected next handler to be called but wasn`t")
	}
}
