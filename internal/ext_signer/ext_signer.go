package ext_signer

import (
	"context"
	"demo-signserver/internal/ext_signer/adapters"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
)

type ExternalSignerAdapter adapters.ExternalSignerAdapter

type ExternalSigner struct {
	adapter ExternalSignerAdapter
}

func NewExternalSigner(ctx context.Context, profileID string) (ExternalSignerAdapter, error) {
	var (
		adapter ExternalSignerAdapter
		err     error
	)
	repo := repositories.NewProfileRepository()
	profile, err := repo.GetProfileByID(profileID)
	if err != nil {
		return nil, err
	}

	switch *profile.Signer {
	case domain.SignerPositivo:
		adapter, err = adapters.NewPositivoSigner(ctx, profile)
		if err != nil {
			return nil, err
		}
	default:
		panic("unsupported signer type")
	}
	return &ExternalSigner{
		adapter: adapter,
	}, nil
}

func (e *ExternalSigner) StartSign(srcPath string) (string, error) {
	return e.adapter.StartSign(srcPath)
}

func (e *ExternalSigner) WaitSignature(ID string, dstPath string) error {
	return e.adapter.WaitSignature(ID, dstPath)
}
