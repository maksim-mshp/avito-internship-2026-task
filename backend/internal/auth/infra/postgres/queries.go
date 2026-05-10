package postgres

const getByRoleQuery = `
	SELECT id::text, email, role, created_at
	FROM ai_assistants_catalog.users
	WHERE id = $1 AND role = $2;
`

const createUserQuery = `
	INSERT INTO ai_assistants_catalog.users (
		id,
		email,
		role,
		password_hash
	)
	VALUES (
		gen_random_uuid(),
		$1,
		$2,
		$3
	)
	RETURNING id::text, email, role, created_at;
`

const getAuthUserByEmailQuery = `
	SELECT id::text, email, role, created_at, password_hash
	FROM ai_assistants_catalog.users
	WHERE email = $1;
`
