package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authenticationv1 "github.com/example/identity-access-service/gen/proto/authentication/v1"
	commonv1 "github.com/example/identity-access-service/gen/proto/common/v1"
	"github.com/example/identity-access-service/internal/authentication/application"
	authdomain "github.com/example/identity-access-service/internal/authentication/domain"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type Controller struct {
	authenticationv1.UnimplementedAuthenticationServiceServer
	validateToken   *application.ValidateTokenUseCase
	getServiceToken *application.GetServiceTokenUseCase
}

func NewController(validateToken *application.ValidateTokenUseCase, getServiceToken *application.GetServiceTokenUseCase) *Controller {
	return &Controller{validateToken: validateToken, getServiceToken: getServiceToken}
}

func (c *Controller) ValidateToken(ctx context.Context, req *authenticationv1.ValidateTokenRequest) (*authenticationv1.ValidateTokenResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	principal, err := c.validateToken.Execute(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	return &authenticationv1.ValidateTokenResponse{Principal: mapPrincipal(principal)}, nil
}

func (c *Controller) GetServiceToken(ctx context.Context, req *authenticationv1.GetServiceTokenRequest) (*authenticationv1.GetServiceTokenResponse, error) {
	if req.GetServiceId() == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}
	scopes := make([]shareddomain.Scope, 0, len(req.GetScopes()))
	for _, scope := range req.GetScopes() {
		scopes = append(scopes, shareddomain.Scope(scope))
	}
	token, err := c.getServiceToken.Execute(ctx, application.GetServiceTokenCommand{ServiceID: req.GetServiceId(), Scopes: scopes})
	if err != nil {
		log.Printf("GER SERVICE TOKEN ERROR %v", err)
		switch {
		case errors.Is(err, authdomain.ErrUnknownService):
			return nil, status.Error(codes.NotFound, "unknown service")
		case errors.Is(err, authdomain.ErrUnAuthorized):
			return nil, status.Error(codes.PermissionDenied, "authorized")
		default:
			return nil, status.Error(codes.Internal, "could not obtain service token")
		}
	}
	return &authenticationv1.GetServiceTokenResponse{AccessToken: token.Value, TokenType: token.TokenType, ExpiresIn: token.ExpiresIn}, nil
}

func mapPrincipal(principal shareddomain.Principal) *commonv1.Principal {
	principalType := commonv1.PrincipalType_UNSPECIFIED
	switch principal.Type {
	case shareddomain.PrincipalTypeUser:
		principalType = commonv1.PrincipalType_USER
	case shareddomain.PrincipalTypeService:
		principalType = commonv1.PrincipalType_SERVICE
	}
	scopes := make([]string, 0, len(principal.Scopes))
	for _, scope := range principal.Scopes {
		scopes = append(scopes, string(scope))
	}
	return &commonv1.Principal{Id: principal.ID, Type: principalType, Scopes: scopes}
}
