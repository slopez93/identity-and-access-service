package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	authdomain "github.com/example/identity-access-service/internal/authorization/domain"
)

type PermissionRepository struct{ db *pgxpool.Pool }

func NewPermissionRepository(db *pgxpool.Pool) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) FindByResourceAndAction(ctx context.Context, resource authdomain.ResourceType, action authdomain.Action) (*authdomain.Permission, error) {
	const query = `
		SELECT id::text, resource, action, requires_ownership
		FROM permissions
		WHERE resource = $1 AND action = $2
		LIMIT 1`

	var permission authdomain.Permission
	err := r.db.QueryRow(ctx, query, resource, action).Scan(&permission.ID, &permission.Resource, &permission.Action, &permission.RequiresOwnership)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, authdomain.ErrPermissionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find permission by resource and action: %w", err)
	}
	return &permission, nil
}
