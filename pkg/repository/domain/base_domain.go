package domain

type BaseDomain struct {
	ID        string `json:"id,omitempty" dynamodbav:"-"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}
