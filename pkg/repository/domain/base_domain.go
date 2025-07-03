package domain

import (
	"time"
)

type BaseDomainInterface interface {
	SetID(string)
	GetID() string
	SetCreateTs()
	SetUpdateTs()
}

type BaseDomain struct {
	ID        string `json:"id" dynamodbav:"-"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt string `json:"updated_at" dynamodbav:"updated_at"`
}

func (b *BaseDomain) GetID() string {
	return b.ID
}

func (b *BaseDomain) SetID(id string) {
	b.ID = id
}

func (b *BaseDomain) SetCreateTs() {
	now := time.Now().UTC().Format(time.RFC3339)
	b.CreatedAt = now
	b.UpdatedAt = now
}

func (b *BaseDomain) SetUpdateTs() {
	now := time.Now().UTC().Format(time.RFC3339)
	b.UpdatedAt = now
}
