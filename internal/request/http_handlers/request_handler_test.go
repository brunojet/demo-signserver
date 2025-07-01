package http_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"demo-signserver/internal/request/dtos"
	"demo-signserver/internal/request/services"
	db_services "demo-signserver/pkg/repository/services"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

var (
	ProfileTable_2 = "profile_handler_2"
	RequestTable_2 = "request_handler_2"
)

func init() {
	gin.SetMode(gin.TestMode)
	r = gin.Default()
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
	os.Setenv("SIGN_PROFILE_TABLE", ProfileTable_2)
	os.Setenv("SIGN_REQUEST_TABLE", RequestTable_2)
	db := db_services.NewDB("Dummy")
	db.DeleteTable(context.TODO(), ProfileTable_2)
	db.DeleteTable(context.TODO(), RequestTable_2)
	db.CreateTable(context.TODO(), ProfileTable_2, "profile_id")
	db.CreateTable(context.TODO(), RequestTable_2, "")
	handler := NewRequestHandler(services.NewRequestService())
	handler.RegisterRoutes(r)
	handler_profile := NewProfileHandler(services.NewProfileService())
	handler_profile.RegisterRoutes(r)
	log.Println("RequestHandler initialized")
}

func TestRequestHandler_CreateRequest_BadRequest(t *testing.T) {
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
	p := httptest.NewRecorder()
	w := httptest.NewRecorder()

	// Cria um profile válido antes
	profilePayload := map[string]interface{}{
		"signer":      "positivo",
		"profile_id":  "456",
		"description": "Profile Teste",
		"upload": map[string]interface{}{
			"url":      "https://upload.com/file",
			"tries":    1,
			"interval": 5,
		},
		"download": map[string]interface{}{
			"url":      "https://download.com/file",
			"tries":    1,
			"interval": 5,
		},
	}
	bProfile, _ := json.Marshal(profilePayload)
	profileReq, _ := http.NewRequest("POST", "/profiles", bytes.NewBuffer(bProfile))
	profileReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(p, profileReq)
	if p.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 for profile, got %d", p.Code)
	}
	var profileResp map[string]interface{}
	_ = json.Unmarshal(p.Body.Bytes(), &profileResp)
	profileID, _ := profileResp["id"].(string)

	// Agora cria a request
	webhook := "https://callback.com/webhook"
	payload := dtos.CreateSignRequestDTO{ProfileId: profileID, WebhookURL: &webhook}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/requests", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
	var resp map[string]interface{}
	log.Println("Response Body:", w.Body.String())
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	if err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Errorf("Expected id in response, got %v", resp)
	}
}

func TestRequestHandler_GetRequestByID_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/requests/inexistente", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
