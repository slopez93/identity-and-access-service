package domain

type PrincipalType string

const (
	PrincipalTypeUser    PrincipalType = "USER"
	PrincipalTypeService PrincipalType = "SERVICE"
)

type Scope string

type Principal struct {
	ID     string
	Type   PrincipalType
	Scopes []Scope
}
