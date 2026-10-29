package interfaces

import (
	"context"

	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type PermissionChecker interface {
	HasPermission(ctx context.Context, principalID string, principalType shareddomain.PrincipalType, permissionID string) (bool, error)
}
