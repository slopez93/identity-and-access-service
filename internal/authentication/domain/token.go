package domain

type AccessToken struct {
	Value     string
	TokenType string
	ExpiresIn int32
}
