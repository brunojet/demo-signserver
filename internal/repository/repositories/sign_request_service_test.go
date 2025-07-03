package repositories

import (
	"context"
	"log"
	"os"
	"testing"

	"demo-signserver/internal/repository/domain"
	db_services "demo-signserver/pkg/repository/services"

	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("SIGN_REQUEST_TABLE", "request")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
	db := db_services.NewDB("Dummy")
	if err := db.DeleteTable(context.Background(), "request"); err != nil {
		log.Fatalf("Error deleting table request: %v", err)
	}
	if err := db.CreateTable(context.Background(), "request", ""); err != nil {
		log.Fatalf("Error creating table request: %v", err)
	}
}

func setupSignRequestServiceTest() *SignRequestService {
	return NewSignRequestService()
}

func TestSignRequestService_CreateRequest_And_GetRequestByID(t *testing.T) {
	svc := setupSignRequestServiceTest()
	request := &domain.SignRequest{}
	err := svc.CreateRequest(request)
	assert.NoError(t, err)

	result, err := svc.GetRequestByID(request.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, request.ID, result.ID)
}

func TestSignRequestService_GetRequestByID_NotFound(t *testing.T) {
	svc := setupSignRequestServiceTest()
	result, err := svc.GetRequestByID("not-exists")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSignRequestService_UpdateRequest(t *testing.T) {
	svc := setupSignRequestServiceTest()
	request := &domain.SignRequest{}
	err := svc.CreateRequest(request)
	assert.NoError(t, err)

	requestUpdated := &domain.SignRequest{}
	err = svc.UpdateRequest(request.ID, requestUpdated)
	assert.NoError(t, err)

	result, err := svc.GetRequestByID(request.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, request.ID, result.ID)
}
