package http_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"demo-signserver/internal/http_handler/dtos"

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

func setupRouter() *gin.Engine {
	r := gin.Default()
	RegisterProfileRoutes(r)
	return r
}

func TestCreateProfile_ValidationError(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{}`) // corpo vazio, deve falhar
	req, _ := http.NewRequest("POST", "/profiles", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProfile_Success(t *testing.T) {
	r := setupRouter()
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
