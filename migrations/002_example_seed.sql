INSERT INTO roles(code) VALUES ('CUSTOMER') ON CONFLICT(code) DO NOTHING;

INSERT INTO role_permissions(role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON (p.resource, p.action) IN (
    ('TRANSACTIONS', 'READ'),
    ('TRANSACTIONS', 'EXECUTE'),
    ('PAYMENTS', 'READ'),
    ('PAYMENTS', 'EXECUTE'),
    ('ACCOUNTS', 'READ')
)
WHERE r.code = 'CUSTOMER'
ON CONFLICT DO NOTHING;

-- Example only. In a real environment principal-role assignments are managed
-- by your administrative workflow.
INSERT INTO principal_roles(principal_id, principal_type, role_id)
SELECT '992e01e4-cf93-4052-ab73-0d87719420e1', 'USER', id FROM roles WHERE code = 'CUSTOMER'
ON CONFLICT DO NOTHING;
