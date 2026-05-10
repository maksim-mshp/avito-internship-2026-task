package postgres

const (
	createPendingQuery = `
		WITH inserted AS (
			INSERT INTO ai_assistants_catalog.assistant_runs (
				assistant_id,
				user_id,
				model,
				user_prompt,
				status
			)
			VALUES ($1, $2, $3, $4, 'pending')
			RETURNING id, assistant_id, user_id, model, user_prompt, output, status, error, created_at
		)
		SELECT
			r.id::text,
			r.assistant_id::text,
			a.name,
			a.category_id::text,
			c.name,
			r.user_id::text,
			r.model,
			r.user_prompt,
			r.output,
			r.status,
			r.error,
			r.created_at
		FROM inserted r
		JOIN ai_assistants_catalog.assistants a ON a.id = r.assistant_id
		JOIN ai_assistants_catalog.categories c ON c.id = a.category_id;
	`

	completeQuery = `
		WITH updated AS (
			UPDATE ai_assistants_catalog.assistant_runs
			SET output = $2, status = 'success', error = NULL
			WHERE id = $1
			RETURNING id, assistant_id, user_id, model, user_prompt, output, status, error, created_at
		)
		SELECT
			r.id::text,
			r.assistant_id::text,
			a.name,
			a.category_id::text,
			c.name,
			r.user_id::text,
			r.model,
			r.user_prompt,
			r.output,
			r.status,
			r.error,
			r.created_at
		FROM updated r
		JOIN ai_assistants_catalog.assistants a ON a.id = r.assistant_id
		JOIN ai_assistants_catalog.categories c ON c.id = a.category_id;
	`

	failQuery = `
		WITH updated AS (
			UPDATE ai_assistants_catalog.assistant_runs
			SET output = NULL, status = 'failed', error = $2
			WHERE id = $1
			RETURNING id, assistant_id, user_id, model, user_prompt, output, status, error, created_at
		)
		SELECT
			r.id::text,
			r.assistant_id::text,
			a.name,
			a.category_id::text,
			c.name,
			r.user_id::text,
			r.model,
			r.user_prompt,
			r.output,
			r.status,
			r.error,
			r.created_at
		FROM updated r
		JOIN ai_assistants_catalog.assistants a ON a.id = r.assistant_id
		JOIN ai_assistants_catalog.categories c ON c.id = a.category_id;
	`

	listMyQuery = `
		SELECT
			r.id::text,
			r.assistant_id::text,
			a.name,
			a.category_id::text,
			c.name,
			r.user_id::text,
			r.model,
			r.user_prompt,
			r.output,
			r.status,
			r.error,
			r.created_at,
			COUNT(*) OVER()
		FROM ai_assistants_catalog.assistant_runs r
		JOIN ai_assistants_catalog.assistants a ON a.id = r.assistant_id
		JOIN ai_assistants_catalog.categories c ON c.id = a.category_id
		WHERE r.user_id = $1
			AND ($2::text IS NULL OR r.status = $2)
		ORDER BY r.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	listAdminQuery = `
		SELECT
			r.id::text,
			r.assistant_id::text,
			a.name,
			a.category_id::text,
			c.name,
			r.user_id::text,
			r.model,
			r.user_prompt,
			r.output,
			r.status,
			r.error,
			r.created_at,
			COUNT(*) OVER()
		FROM ai_assistants_catalog.assistant_runs r
		JOIN ai_assistants_catalog.assistants a ON a.id = r.assistant_id
		JOIN ai_assistants_catalog.categories c ON c.id = a.category_id
		WHERE ($1::uuid IS NULL OR r.assistant_id = $1)
			AND ($2::text IS NULL OR r.status = $2)
		ORDER BY r.created_at DESC
		LIMIT $3 OFFSET $4;
	`
)
