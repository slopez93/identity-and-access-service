package application

import (
	"context"
	"testing"

	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type verifierStub struct{ principal shareddomain.Principal }

func (v verifierStub) Verify(context.Context, string) (shareddomain.Principal, error) {
	return v.principal, nil
}

func TestValidateTokenDelegatesToVerifier(t *testing.T) {
	expected := shareddomain.Principal{ID: "user-123", Type: shareddomain.PrincipalTypeUser}
	uc := NewValidateTokenUseCase(verifierStub{principal: expected})
	got, err := uc.Execute(context.Background(), "token")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != expected.ID || got.Type != expected.Type {
		t.Fatalf("unexpected principal: %+v", got)
	}
}
