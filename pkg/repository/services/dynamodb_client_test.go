package db_services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDynamoDBClient_Success(t *testing.T) {
	client := GetDynamoDBCLient()
	assert.NotNil(t, client)
}
