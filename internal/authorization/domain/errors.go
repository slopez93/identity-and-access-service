package domain

import "errors"

var (
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrScopeMappingMissing = errors.New("required scope mapping not found")
)
