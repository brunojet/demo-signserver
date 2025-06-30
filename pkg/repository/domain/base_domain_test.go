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

func TestBaseDomain_SetCreateTs(t *testing.T) {
	b := &BaseDomain{}
	b.SetCreateTs()
	if b.CreatedAt == 0 || b.UpdatedAt == 0 {
		t.Errorf("SetCreateTs não preencheu os campos corretamente: CreatedAt=%d, UpdatedAt=%d", b.CreatedAt, b.UpdatedAt)
	}
	if b.CreatedAt != b.UpdatedAt {
		t.Errorf("CreatedAt e UpdatedAt deveriam ser iguais após SetCreateTs, mas são diferentes: %d vs %d", b.CreatedAt, b.UpdatedAt)
	}
}

func TestBaseDomain_SetUpdateTs(t *testing.T) {
	b := &BaseDomain{}
	b.SetCreateTs()
	oldCreated := b.CreatedAt
	oldUpdated := b.UpdatedAt
	b.SetUpdateTs()
	if b.UpdatedAt == 0 {
		t.Errorf("SetUpdateTs não preencheu UpdatedAt corretamente")
	}
	if b.UpdatedAt < oldUpdated {
		t.Errorf("UpdatedAt deveria ser maior ou igual ao valor anterior")
	}
	if b.CreatedAt != oldCreated {
		t.Errorf("CreatedAt não deveria ser alterado por SetUpdateTs")
	}
}

func TestBaseDomainInterface_IsImplemented(t *testing.T) {
	var _ BaseDomainInterface = (*BaseDomain)(nil)
}
