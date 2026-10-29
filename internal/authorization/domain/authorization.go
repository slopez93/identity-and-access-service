package domain

import shareddomain "github.com/example/identity-access-service/internal/shared/domain"

type DenyReason string

const (
	DenyMissingScope        DenyReason = "MISSING_SCOPE"
	DenyMissingPermission   DenyReason = "MISSING_PERMISSION"
	DenyPrincipalNotAllowed DenyReason = "PRINCIPAL_NOT_ALLOWED"
)

type AuthorizationRequest struct {
	Principal       shareddomain.Principal
	Resource        ResourceType
	Action          Action
	ResourceID      string
	ResourceOwnerID string
}

type AuthorizationResult struct {
	Allowed bool
	Reason  DenyReason
}
