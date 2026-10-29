package application

import (
	"context"
	"fmt"

	authdomain "github.com/example/identity-access-service/internal/authorization/domain"
	"github.com/example/identity-access-service/internal/authorization/domain/interfaces"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type TransactionAction int

const (
	TransactionActionUnspecified TransactionAction = iota
	TransactionActionRead
	TransactionActionExecute
	TransactionActionCreate
	TransactionActionDelete
)

type TransactionResource struct {
	ID      string
	OwnerID string
}

type AuthorizeTransactionCommand struct {
	Principal shareddomain.Principal
	Action    TransactionAction
	Resource  TransactionResource
}

type AuthorizeTransactionUseCase struct {
	scopes             interfaces.RequiredScopeProvider
	permissions        interfaces.PermissionRepository
	permissionsChecker interfaces.PermissionChecker
}

func NewAuthorizeTransactionUseCase(scopes interfaces.RequiredScopeProvider, permissions interfaces.PermissionRepository, permissionsChecker interfaces.PermissionChecker) *AuthorizeTransactionUseCase {
	return &AuthorizeTransactionUseCase{scopes: scopes, permissions: permissions, permissionsChecker: permissionsChecker}
}

func (uc *AuthorizeTransactionUseCase) Execute(ctx context.Context, cmd AuthorizeTransactionCommand) (authdomain.AuthorizationResult, error) {
	action, err := mapTransactionAction(cmd.Action)
	if err != nil {
		return authdomain.AuthorizationResult{}, err
	}

	permission, err := uc.permissions.FindByResourceAndAction(ctx, authdomain.ResourceType(cmd.Resource.ID), action)
	if err != nil {
		return authdomain.AuthorizationResult{}, err
	}

	requiredScopes, err := uc.scopes.ForPermission(*permission)
	if err != nil {
		return authdomain.AuthorizationResult{}, err
	}

	if !uc.hasRequiredScopes(cmd.Principal.Scopes, requiredScopes) {
		return authdomain.AuthorizationResult{Allowed: false, Reason: authdomain.DenyMissingScope}, nil
	}

	hasPermission, err := uc.permissionsChecker.HasPermission(ctx, cmd.Principal.ID, cmd.Principal.Type, permission.ID)
	if err != nil {
		return authdomain.AuthorizationResult{}, nil
	}

	if !hasPermission {
		return authdomain.AuthorizationResult{Allowed: false, Reason: authdomain.DenyMissingPermission}, nil
	}

	if permission.RequiresOwnership && cmd.Principal.ID != cmd.Resource.OwnerID {
		return authdomain.AuthorizationResult{Allowed: false, Reason: authdomain.DenyPrincipalNotAllowed}, nil
	}

	return authdomain.AuthorizationResult{Allowed: true}, nil

}

func (uc *AuthorizeTransactionUseCase) hasRequiredScopes(actual []shareddomain.Scope, required []shareddomain.Scope) bool {
	set := make(map[shareddomain.Scope]bool, len(actual))

	for _, scope := range actual {
		set[scope] = true
	}

	for _, scope := range required {
		if _, ok := set[scope]; !ok {
			return false
		}
	}

	return true
}

func mapTransactionAction(action TransactionAction) (authdomain.Action, error) {
	switch action {
	case TransactionActionRead:
		return authdomain.ActionRead, nil
	case TransactionActionCreate:
		return authdomain.ActionCreate, nil
	case TransactionActionExecute:
		return authdomain.ActionExecute, nil
	default:
		return "", fmt.Errorf("unsupported transaction action: %d", action)
	}
}
