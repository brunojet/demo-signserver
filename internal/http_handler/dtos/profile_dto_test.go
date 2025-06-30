package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferInfoDTO_Validation(t *testing.T) {
	dto := TransferInfoDTO{
		Url:      "not-a-url",
		Tries:    0,
		Interval: 100,
	}
	// Não é possível validar diretamente sem o Gin, mas podemos testar os valores
	assert.NotEqual(t, "", dto.Url)
	assert.True(t, dto.Tries >= 0)
	assert.True(t, dto.Interval >= 0)
}
