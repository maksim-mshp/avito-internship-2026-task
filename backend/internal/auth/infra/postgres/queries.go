package postgres

const getByRoleQuery = `
	SELECT id::text, email, role, created_at
	FROM ai_assistants_catalog.users
	WHERE id = $1 AND role = $2;
`
