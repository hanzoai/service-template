// Package migrations holds DB schema migrations applied per-(org,user)
// on first hydrate. The skeleton ships one placeholder migration that
// creates the Example table.
package migrations

// InitialSchema is the first migration. Each service maintains its own
// migrations; the base migration runner stamps each SQLite with the
// schema version on apply.
const InitialSchema = `
CREATE TABLE IF NOT EXISTS examples (
	id         TEXT PRIMARY KEY,
	name       TEXT NOT NULL,
	created_at INTEGER NOT NULL
);
`
