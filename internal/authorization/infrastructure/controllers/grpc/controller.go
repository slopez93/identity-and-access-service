package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authorizationv1 "github.com/example/identity-access-service/gen/proto/authorization/v1"
	commonv1 "github.com/example/identity-access-service/gen/proto/common/v1"
	"github.com/example/identity-access-service/internal/authorization/application"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type Controller struct {
	authorizationv1.UnimplementedAuthorizationServiceServer
	authorizeTransaction *application.AuthorizeTransactionUseCase
}

func NewController(authorizeTransaction *application.AuthorizeTransactionUseCase) *Controller {
	return &Controller{authorizeTransaction: authorizeTransaction}
}

func (c *Controller) AuthorizeTransaction(ctx context.Context, req *authorizationv1.AuthorizeTransactionRequest) (*authorizationv1.AuthorizationResponse, error) {
	cmd, err := mapAuthorizeTransactionRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := c.authorizeTransaction.Execute(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Internal, "authorization failed")
	}
	return &authorizationv1.AuthorizationResponse{Allowed: result.Allowed, Reason: string(result.Reason)}, nil
}

func mapAuthorizeTransactionRequest(req *authorizationv1.AuthorizeTransactionRequest) (application.AuthorizeTransactionCommand, error) {
	if req.GetPrincipal() == nil || req.GetResource() == nil {
		return application.AuthorizeTransactionCommand{}, status.Error(codes.InvalidArgument, "principal and resource are required")
	}
	principal, err := mapPrincipal(req.GetPrincipal())
	if err != nil {
		return application.AuthorizeTransactionCommand{}, err
	}
	action, err := mapTransactionAction(req.GetAction())
	if err != nil {
		return application.AuthorizeTransactionCommand{}, err
	}
	return application.AuthorizeTransactionCommand{
		Principal: principal,
		Action:    action,
		Resource:  application.TransactionResource{ID: req.GetResource().GetId(), OwnerID: req.GetResource().GetOwnerId()},
	}, nil
}

func mapPrincipal(p *commonv1.Principal) (shareddomain.Principal, error) {
	principalType := shareddomain.PrincipalType("")
	switch p.GetType() {
	case commonv1.PrincipalType_USER:
		principalType = shareddomain.PrincipalTypeUser
	case commonv1.PrincipalType_SERVICE:
		principalType = shareddomain.PrincipalTypeService
	default:
		return shareddomain.Principal{}, status.Error(codes.InvalidArgument, "invalid principal type")
	}
	scopes := make([]shareddomain.Scope, 0, len(p.GetScopes()))
	for _, scope := range p.GetScopes() {
		scopes = append(scopes, shareddomain.Scope(scope))
	}
	return shareddomain.Principal{ID: p.GetId(), Type: principalType, Scopes: scopes}, nil
}

func mapTransactionAction(action authorizationv1.TransactionAction) (application.TransactionAction, error) {
	switch action {
	case authorizationv1.TransactionAction_READ:
		return application.TransactionActionRead, nil
	case authorizationv1.TransactionAction_EXECUTE:
		return application.TransactionActionExecute, nil
	case authorizationv1.TransactionAction_CREATE:
		return application.TransactionActionCreate, nil
	default:
		return application.TransactionActionUnspecified, status.Error(codes.InvalidArgument, "unsupported transaction action")
	}
}
