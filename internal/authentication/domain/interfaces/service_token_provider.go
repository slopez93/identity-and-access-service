package interfaces

import (
	"context"

	authdomain "github.com/example/identity-access-service/internal/authentication/domain"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type ServiceTokenProvider interface {
	GetToken(ctx context.Context, serviceID string, scopes []shareddomain.Scope) (authdomain.AccessToken, error)
}
