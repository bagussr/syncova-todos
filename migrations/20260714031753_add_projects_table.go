package migrations

import "gorm.io/gorm"

var Migration20260714031753 = Migration{
	ID:     "20260714031753",
	Before: "",
	Upgrade: func(db *gorm.DB) error {
		db.Exec(
			"CREATE TYPE status AS ENUM ('not_started', 'in_progress', 'testing', 'completed', 'backlog')",
		)
		db.Exec(
			"CREATE TABLE projects (" +
				"id SERIAL PRIMARY KEY," +
				"uuid UUID DEFAULT gen_random_uuid()," +
				"status status DEFAULT 'not_started' NOT NULL," +
				"due_date DATE" +
				")",
		)
		return nil
	},
	Downgrade: func(db *gorm.DB) error {
		db.Exec("DROP TABLE IF EXISTS projects")
		db.Exec("DROP TYPE IF EXISTS status")
		return nil
	},
}

func init() {
	Register(Migration20260714031753)
}
