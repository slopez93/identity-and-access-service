package interfaces

import (
	"context"

	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (shareddomain.Principal, error)
}
