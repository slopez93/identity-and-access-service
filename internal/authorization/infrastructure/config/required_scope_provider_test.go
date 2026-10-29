package config

import (
	"testing"

	authdomain "github.com/example/identity-access-service/internal/authorization/domain"
	sharedconfig "github.com/example/identity-access-service/internal/shared/config"
)

func TestRequiredScopeProvider(t *testing.T) {
	provider := NewRequiredScopeProvider(sharedconfig.AuthorizationScopesConfig{TransactionExecute: "financial-api/transactions.execute"})
	scopes, err := provider.ForPermission(authdomain.Permission{Resource: authdomain.ResourceTransaction, Action: authdomain.ActionExecute})
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 1 || scopes[0] != "financial-api/transactions.execute" {
		t.Fatalf("unexpected scopes: %+v", scopes)
	}
}
