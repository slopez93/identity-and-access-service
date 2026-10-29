package interfaces

import (
	authdomain "github.com/example/identity-access-service/internal/authorization/domain"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type RequiredScopeProvider interface {
	ForPermission(permission authdomain.Permission) ([]shareddomain.Scope, error)
}
