package observability

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func resetMetrics() {
	metrics := GetMetricsSnapshot()
	for k := range metrics {
		requestCount[k] = 0
	}
}

func TestObservableMiddleware_RequestIDGeneration(t *testing.T) {
	resetMetrics()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var capturedRequestID string
	r.GET("/test", ObservableMiddleware(func(c *gin.Context, obs *Observability) {
		capturedRequestID = obs.RequestID
		c.String(200, "ok")
	}))

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, capturedRequestID)
	assert.Equal(t, w.Header().Get("X-Request-ID"), capturedRequestID)
}

func TestObservableMiddleware_UsesProvidedRequestID(t *testing.T) {
	resetMetrics()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var capturedRequestID string
	r.GET("/test", ObservableMiddleware(func(c *gin.Context, obs *Observability) {
		capturedRequestID = obs.RequestID
		c.String(200, "ok")
	}))

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "my-custom-id")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "my-custom-id", capturedRequestID)
	assert.Equal(t, "my-custom-id", w.Header().Get("X-Request-ID"))
}

func TestObservableMiddleware_LogsAndMetrics(t *testing.T) {
	resetMetrics()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var buf bytes.Buffer
	defaultLogger.SetOutput(&buf)
	defer defaultLogger.SetOutput(os.Stdout)

	r.GET("/test", ObservableMiddleware(func(c *gin.Context, obs *Observability) {
		c.String(200, "ok")
	}))

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Aguarda logs assíncronos
	time.Sleep(10 * time.Millisecond)

	// Verifica se logou entrada e saída
	logs := buf.String()
	assert.Contains(t, logs, "Request recebida")
	assert.Contains(t, logs, "Request finalizada")
	// Verifica se incrementou métrica
	metrics := GetMetricsSnapshot()
	assert.True(t, metrics["/test"] > 0)
}

func TestObservableMiddleware_HandlerIsCalled(t *testing.T) {
	resetMetrics()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	called := false
	r.GET("/test", ObservableMiddleware(func(c *gin.Context, obs *Observability) {
		called = true
		c.String(200, "ok")
	}))

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, 200, w.Code)
}
