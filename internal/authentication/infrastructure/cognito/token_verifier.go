package cognito

import (
	"context"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"

	authdomain "github.com/example/identity-access-service/internal/authentication/domain"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type TokenVerifierConfig struct {
	Issuer           string
	UserClientID     string
	ServiceClientIDs map[string]string // Cognito client_id -> internal service_id
}

type cognitoClaims struct {
	jwt.RegisteredClaims
	ClientID string `json:"client_id"`
	TokenUse string `json:"token_use"`
	Scope    string `json:"scope"`
}

type TokenVerifier struct {
	issuer           string
	jwks             keyfunc.Keyfunc
	userClientID     string
	serviceClientIDs map[string]string
}

func NewTokenVerifier(ctx context.Context, cfg TokenVerifierConfig) (*TokenVerifier, error) {
	jwksURL := strings.TrimRight(cfg.Issuer, "/") + "/.well-known/jwks.json"
	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, err
	}
	return &TokenVerifier{
		issuer:           cfg.Issuer,
		jwks:             jwks,
		userClientID:     cfg.UserClientID,
		serviceClientIDs: cfg.ServiceClientIDs,
	}, nil
}

func (v *TokenVerifier) Verify(ctx context.Context, rawToken string) (shareddomain.Principal, error) {
	claims := &cognitoClaims{}
	token, err := jwt.ParseWithClaims(
		rawToken,
		claims,
		v.jwks.KeyfuncCtx(ctx),
		jwt.WithIssuer(v.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil || !token.Valid {
		return shareddomain.Principal{}, authdomain.ErrInvalidToken
	}
	if claims.TokenUse != "access" {
		return shareddomain.Principal{}, authdomain.ErrInvalidTokenUse
	}

	scopes := parseScopes(claims.Scope)
	if serviceID, ok := v.serviceClientIDs[claims.ClientID]; ok {
		return shareddomain.Principal{ID: serviceID, Type: shareddomain.PrincipalTypeService, Scopes: scopes}, nil
	}
	if claims.ClientID != v.userClientID {
		return shareddomain.Principal{}, authdomain.ErrInvalidClient
	}
	if claims.Subject == "" {
		return shareddomain.Principal{}, authdomain.ErrInvalidToken
	}
	return shareddomain.Principal{ID: claims.Subject, Type: shareddomain.PrincipalTypeUser, Scopes: scopes}, nil
}

func parseScopes(raw string) []shareddomain.Scope {
	values := strings.Fields(raw)
	result := make([]shareddomain.Scope, 0, len(values))
	for _, value := range values {
		result = append(result, shareddomain.Scope(value))
	}
	return result
}
