package http_handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"demo-signserver/internal/request/dtos"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("SIGN_PROFILE_TABLE", "profile")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
	os.Setenv("DELETE_TABLE", "true")
}

func SetupProfileRouter() *gin.Engine {
	r := gin.Default()
	RegisterProfileRoutes(r)
	return r
}

func TestCreateProfile_ValidationError(t *testing.T) {
	r := SetupProfileRouter()
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{}`) // corpo vazio, deve falhar
	req, _ := http.NewRequest("POST", "/profiles", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProfile_Success(t *testing.T) {
	r := SetupProfileRouter()
	w := httptest.NewRecorder()
	profile := dtos.CreateSignerProfileDTO{
		Signer:      "positivo",
		ProfileId:   "123",
		Description: "Perfil de teste",
		Upload:      dtos.TransferInfoDTO{Url: "http://teste.com", Tries: 1, Interval: 5},
		Download:    dtos.TransferInfoDTO{Url: "http://teste.com", Tries: 1, Interval: 5},
	}
	jsonBody, _ := json.Marshal(profile)
	req, _ := http.NewRequest("POST", "/profiles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "id")
}

func TestGetProfileByID_NotFound(t *testing.T) {
	r := SetupProfileRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profiles/999", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Profile not found")
}

func TestUpdateProfile_ValidationError(t *testing.T) {
	r := SetupProfileRouter()
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{}`) // corpo vazio, mas PATCH permite update parcial, então deve aceitar
	req, _ := http.NewRequest("PATCH", "/profiles/123", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	// Como não existe o perfil, deve retornar erro interno ou not found
	assert.True(t, w.Code == http.StatusInternalServerError || w.Code == http.StatusBadRequest)
}

func TestUpdateProfile_Success(t *testing.T) {
	r := SetupProfileRouter()
	w := httptest.NewRecorder()
	// Primeiro cria o perfil
	profile := dtos.CreateSignerProfileDTO{
		Signer:      "positivo",
		ProfileId:   "124",
		Description: "Perfil para update",
		Upload:      dtos.TransferInfoDTO{Url: "http://teste.com", Tries: 1, Interval: 5},
		Download:    dtos.TransferInfoDTO{Url: "http://teste.com", Tries: 1, Interval: 5},
	}
	jsonBody, _ := json.Marshal(profile)
	req, _ := http.NewRequest("POST", "/profiles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Extrai o id da resposta
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	id := resp["id"].(string)

	// Agora faz o update
	w2 := httptest.NewRecorder()
	update := dtos.UpdateSignerProfileDTO{
		Description: ptrString("Novo nome"),
	}
	jsonUpdate, _ := json.Marshal(update)
	req2, _ := http.NewRequest("PATCH", "/profiles/"+url.PathEscape(id), bytes.NewBuffer(jsonUpdate))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Profile updated")
}

func ptrString(s string) *string { return &s }
