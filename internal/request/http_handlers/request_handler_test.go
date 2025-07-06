package http_handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"demo-signserver/internal/request/dtos"
	db_services "demo-signserver/pkg/repository/services"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, http.StatusCreated, p.Code, "Expected status 201 for profile creation")
	var profileResp map[string]interface{}
	err := json.Unmarshal(p.Body.Bytes(), &profileResp)
	assert.NoError(t, err, "Error parsing profile response body")
	profileID, _ := profileResp["id"].(string)

	// Agora cria a request
	webhook := "https://callback.com/webhook"
	payload := dtos.CreateSignRequestDTO{ProfileId: profileID, WebhookURL: &webhook}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/requests", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 for request creation")
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err, "Error parsing response body")
	assert.NotEmpty(t, resp[db_services.ID_KEY], "Expected id in response")
}

func TestRequestHandler_GetRequestByID_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/requests/inexistente", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status 404 for non-existent request")
}
