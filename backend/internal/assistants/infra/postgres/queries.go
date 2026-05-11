package postgres

const (
	listQuery = `
		SELECT
			a.id::text,
			a.category_id::text,
			c.name,
			a.name,
			a.description,
			a.model,
			a.system_prompt,
			a.example_user_prompt,
			a.tags,
			a.is_active,
			a.created_at,
			a.updated_at,
			COUNT(*) OVER()
		FROM ai_assistants_catalog.assistants a
		JOIN ai_assistants_catalog.categories c ON c.id = a.category_id
		WHERE ($1::uuid IS NULL OR a.category_id = $1)
			AND ($2::text IS NULL OR a.name ILIKE '%' || $2 || '%' OR a.description ILIKE '%' || $2 || '%' OR EXISTS (
				SELECT 1
				FROM unnest(a.tags) tag
				WHERE tag ILIKE '%' || $2 || '%'
			))
			AND ($3::text IS NULL OR EXISTS (
				SELECT 1
				FROM unnest(a.tags) tag
				WHERE LOWER(tag) = LOWER($3)
			))
			AND ($4::boolean OR a.is_active)
		ORDER BY a.created_at DESC, a.name ASC
		LIMIT $5 OFFSET $6;
	`

	getByIDQuery = `
		SELECT
			a.id::text,
			a.category_id::text,
			c.name,
			a.name,
			a.description,
			a.model,
			a.system_prompt,
			a.example_user_prompt,
			a.tags,
			a.is_active,
			a.created_at,
			a.updated_at
		FROM ai_assistants_catalog.assistants a
		JOIN ai_assistants_catalog.categories c ON c.id = a.category_id
		WHERE a.id = $1
			AND ($2::boolean OR a.is_active);
	`

	createQuery = `
		INSERT INTO ai_assistants_catalog.assistants (
			category_id,
			name,
			description,
			model,
			system_prompt,
			example_user_prompt,
			tags,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id::text,
			category_id::text,
			NULL::text,
			name,
			description,
			model,
			system_prompt,
			example_user_prompt,
			tags,
			is_active,
			created_at,
			updated_at;
	`

	updateQuery = `
		UPDATE ai_assistants_catalog.assistants
		SET
			category_id = $2,
			name = $3,
			description = $4,
			model = $5,
			system_prompt = $6,
			example_user_prompt = $7,
			tags = $8,
			is_active = $9,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id::text,
			category_id::text,
			NULL::text,
			name,
			description,
			model,
			system_prompt,
			example_user_prompt,
			tags,
			is_active,
			created_at,
			updated_at;
	`
)
