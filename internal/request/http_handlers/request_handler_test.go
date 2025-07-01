package http_handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"demo-signserver/internal/request/dtos"
	"demo-signserver/internal/request/services"

	"github.com/gin-gonic/gin"
)

func SetupRequestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := NewRequestHandler(services.NewRequestService())
	handler.RegisterRoutes(r)
	return r
}

func TestRequestHandler_CreateRequest_BadRequest(t *testing.T) {
	r := SetupRequestRouter()
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"profile_id": 123}`) // profile_id deve ser string
	req, _ := http.NewRequest("POST", "/requests", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRequestHandler_CreateRequest_Success(t *testing.T) {
	r := SetupRequestRouter()
	w := httptest.NewRecorder()
	webhook := "https://callback.com/webhook"
	payload := dtos.CreateSignRequestDTO{ProfileId: "123", WebhookURL: &webhook}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/requests", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["id"] == nil || resp["id"] == "" {
		t.Errorf("Expected id in response, got %v", resp)
	}
}

func TestRequestHandler_GetRequestByID_NotFound(t *testing.T) {
	r := SetupRequestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/requests/inexistente", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
