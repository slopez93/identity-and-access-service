package domain

import "errors"

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidTokenUse      = errors.New("invalid token use")
	ErrInvalidClient        = errors.New("invalid client")
	ErrUnknownService       = errors.New("unknown service")
	ErrUnAuthorized         = errors.New("unauthourized")
	ErrScopeNotAllowed      = errors.New("scope not allowed")
	ErrInvalidTokenResponse = errors.New("invalid token response")
)
