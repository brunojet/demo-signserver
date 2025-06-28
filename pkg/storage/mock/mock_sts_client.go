// MockSTSClient é um mock do client STS para uso em testes.
// Permite simular respostas customizadas para AssumeRole.
// Exemplo de uso:
//
//	mock := &MockSTSClient{
//	  AssumeRoleFunc: func(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
//	    // comportamento customizado
//	  },
//	}
package storage_mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type MockSTSClient struct {
	// AssumeRoleFunc permite customizar o comportamento do mock em cada teste.
	AssumeRoleFunc func(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

// AssumeRole implementa a interface do client STS, usando o comportamento customizado se definido.
func (m *MockSTSClient) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	if m.AssumeRoleFunc != nil {
		return m.AssumeRoleFunc(ctx, params, optFns...)
	}
	return &sts.AssumeRoleOutput{
		Credentials: &types.Credentials{
			AccessKeyId:     awsString("mock-access-key"),
			SecretAccessKey: awsString("mock-secret-key"),
			SessionToken:    awsString("mock-session-token"),
		},
	}, nil
}

// awsString é uma função utilitária para criar ponteiros de string.
func awsString(s string) *string { return &s }
