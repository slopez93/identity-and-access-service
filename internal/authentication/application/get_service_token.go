package application

import (
	"context"

	authdomain "github.com/example/identity-access-service/internal/authentication/domain"
	authinterfaces "github.com/example/identity-access-service/internal/authentication/domain/interfaces"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type GetServiceTokenCommand struct {
	ServiceID string
	Scopes    []shareddomain.Scope
}

type GetServiceTokenUseCase struct {
	tokenProvider authinterfaces.ServiceTokenProvider
}

func NewGetServiceTokenUseCase(tokenProvider authinterfaces.ServiceTokenProvider) *GetServiceTokenUseCase {
	return &GetServiceTokenUseCase{tokenProvider: tokenProvider}
}

func (uc *GetServiceTokenUseCase) Execute(ctx context.Context, cmd GetServiceTokenCommand) (authdomain.AccessToken, error) {
	return uc.tokenProvider.GetToken(ctx, cmd.ServiceID, cmd.Scopes)
}
