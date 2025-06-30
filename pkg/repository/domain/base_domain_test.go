package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseDomain_SetID(t *testing.T) {
	b := &BaseDomain{}
	b.SetID("abc123")
	assert.Equal(t, "abc123", b.ID)
}

func TestBaseDomainInterface_IsImplemented(t *testing.T) {
	var _ BaseDomainInterface = (*BaseDomain)(nil)
}
