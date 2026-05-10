package postgres

const (
	getAllQuery = `
		SELECT id::text, name, description, created_at
		FROM ai_assistants_catalog.categories
		ORDER BY created_at DESC, name ASC;
	`

	createQuery = `
		INSERT INTO ai_assistants_catalog.categories (name, description)
		VALUES ($1, $2)
		RETURNING id::text, name, description, created_at;
	`
)
