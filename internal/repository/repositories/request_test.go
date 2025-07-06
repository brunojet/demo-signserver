package repositories

import (
	"testing"

	"demo-signserver/internal/repository/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupSignRequestServiceTest() *RequestRepository {
	return NewRequestRepository()
}

func TestSignRequestService_CreateRequest_And_GetRequestByID(t *testing.T) {
	svc := setupSignRequestServiceTest()
	request := &domain.SignRequest{}
	request.SetID(uuid.New().String())
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
	request.SetID(uuid.New().String())
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
