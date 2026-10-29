package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	shareddomain "github.com/example/identity-access-service/internal/shared/domain"
)

type PermissionChecker struct{ db *pgxpool.Pool }

func NewPermissionChecker(db *pgxpool.Pool) *PermissionChecker {
	return &PermissionChecker{db: db}
}

func (r *PermissionChecker) HasPermission(ctx context.Context, principalID string, principalType shareddomain.PrincipalType, permissionID string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM principal_roles pr
			JOIN role_permissions rp ON rp.role_id = pr.role_id
			WHERE pr.principal_id = $1
			  AND pr.principal_type = $2
			  AND rp.permission_id = $3::uuid
		)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, principalID, principalType, permissionID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check principal permission: %w", err)
	}
	return exists, nil
}
