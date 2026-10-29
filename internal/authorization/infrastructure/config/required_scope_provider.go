package config

import (
	authdomain "github.com/example/identity-access-service/internal/authorization/domain"
	sharedconfig "github.com/example/identity-access-service/internal/shared/config"
	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type permissionKey struct {
	Resource authdomain.ResourceType
	Action   authdomain.Action
}

type RequiredScopeProvider struct {
	mappings map[permissionKey][]shareddomain.Scope
}

func NewRequiredScopeProvider(cfg sharedconfig.AuthorizationScopesConfig) *RequiredScopeProvider {
	return &RequiredScopeProvider{mappings: map[permissionKey][]shareddomain.Scope{
		{authdomain.ResourceTransaction, authdomain.ActionRead}:    {shareddomain.Scope(cfg.TransactionRead)},
		{authdomain.ResourceTransaction, authdomain.ActionExecute}: {shareddomain.Scope(cfg.TransactionExecute)},
		{authdomain.ResourcePayment, authdomain.ActionRead}:        {shareddomain.Scope(cfg.PaymentRead)},
		{authdomain.ResourcePayment, authdomain.ActionExecute}:     {shareddomain.Scope(cfg.PaymentExecute)},
		{authdomain.ResourceAccount, authdomain.ActionRead}:        {shareddomain.Scope(cfg.AccountRead)},
		{authdomain.ResourceAccount, authdomain.ActionUpdate}:      {shareddomain.Scope(cfg.AccountUpdate)},
	}}
}

func (p *RequiredScopeProvider) ForPermission(permission authdomain.Permission) ([]shareddomain.Scope, error) {
	scopes, ok := p.mappings[permissionKey{Resource: permission.Resource, Action: permission.Action}]
	if !ok {
		return nil, authdomain.ErrScopeMappingMissing
	}
	return scopes, nil
}
