package interfaces

import (
	"context"

	authdomain "github.com/example/identity-access-service/internal/authorization/domain"
)

type PermissionRepository interface {
	FindByResourceAndAction(ctx context.Context, resource authdomain.ResourceType, action authdomain.Action) (*authdomain.Permission, error)
}
