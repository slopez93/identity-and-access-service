package application

import (
	"context"

	authinterfaces "github.com/example/identity-access-service/internal/authentication/domain/interfaces"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type ValidateTokenUseCase struct {
	tokenVerifier authinterfaces.TokenVerifier
}

func NewValidateTokenUseCase(tokenVerifier authinterfaces.TokenVerifier) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{tokenVerifier: tokenVerifier}
}

func (uc *ValidateTokenUseCase) Execute(ctx context.Context, token string) (shareddomain.Principal, error) {
	return uc.tokenVerifier.Verify(ctx, token)
}
