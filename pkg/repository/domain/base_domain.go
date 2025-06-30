package domain

import "time"

type BaseDomainInterface interface {
	SetID(string)
	SetCreateTs()
	SetUpdateTs()
}

type BaseDomain struct {
	ID        string `json:"id,omitempty" dynamodbav:"-"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

func (b *BaseDomain) SetID(id string) {
	b.ID = id
}

func (b *BaseDomain) SetCreateTs() {
	now := time.Now().Unix()
	b.CreatedAt = now
	b.UpdatedAt = now
}

func (b *BaseDomain) SetUpdateTs() {
	now := time.Now().Unix()
	b.UpdatedAt = now
}
