CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    requires_ownership BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(resource, action)
);

CREATE TABLE IF NOT EXISTS principal_roles (
    principal_id VARCHAR(255) NOT NULL,
    principal_type VARCHAR(20) NOT NULL,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY(principal_id, principal_type, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY(role_id, permission_id)
);

CREATE INDEX IF NOT EXISTS idx_principal_roles_principal
    ON principal_roles(principal_id, principal_type);

CREATE INDEX IF NOT EXISTS idx_role_permissions_permission
    ON role_permissions(permission_id);

INSERT INTO permissions(resource, action, requires_ownership) VALUES
    ('TRANSACTIONS', 'READ', TRUE),
    ('TRANSACTIONS', 'EXECUTE', TRUE),
    ('TRANSACTIONS', 'APPROVE', FALSE),
    ('PAYMENTS', 'READ', TRUE),
    ('PAYMENTS', 'EXECUTE', TRUE),
    ('ACCOUNTS', 'READ', TRUE),
    ('ACCOUNTS', 'UPDATE', TRUE)
ON CONFLICT(resource, action) DO NOTHING;
