package cognito

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	authdomain "github.com/example/identity-access-service/internal/authentication/domain"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type Client interface {
	GetClientToken(ctx context.Context, params *cognitoidentityprovider.GetClientTokenInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GetClientTokenOutput, error)
}

type ServiceClient struct {
	ClientID     string
	ClientSecret string
}

type ServiceTokenProvider struct {
	client   Client
	services map[string]ServiceClient
}

func NewServiceTokenProvider(client Client, services map[string]ServiceClient) *ServiceTokenProvider {
	return &ServiceTokenProvider{client: client, services: services}
}

func (p *ServiceTokenProvider) GetToken(ctx context.Context, serviceID string, scopes []shareddomain.Scope) (authdomain.AccessToken, error) {
	service, ok := p.services[serviceID]
	if !ok {
		return authdomain.AccessToken{}, authdomain.ErrUnknownService
	}

	requestedScopes := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		requestedScopes = append(requestedScopes, string(scope))
	}

	out, err := p.client.GetClientToken(ctx, &cognitoidentityprovider.GetClientTokenInput{
		ClientId: aws.String(service.ClientID),
		Secret:   aws.String(service.ClientSecret),
		Scopes:   requestedScopes,
	})
	if err != nil {
		var notAuthorized *types.NotAuthorizedException

		if errors.As(err, &notAuthorized) {
			return authdomain.AccessToken{}, authdomain.ErrUnAuthorized
		}

		return authdomain.AccessToken{}, fmt.Errorf("get Cognito client token: %w", err)
	}
	if out.ClientAuthenticationResult == nil || out.ClientAuthenticationResult.AccessToken == nil {
		return authdomain.AccessToken{}, authdomain.ErrInvalidTokenResponse
	}

	return authdomain.AccessToken{
		Value:     aws.ToString(out.ClientAuthenticationResult.AccessToken),
		TokenType: aws.ToString(out.ClientAuthenticationResult.TokenType),
		ExpiresIn: out.ClientAuthenticationResult.ExpiresIn,
	}, nil
}
