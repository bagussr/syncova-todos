package migrations

import "gorm.io/gorm"

var Migration20260715021827 = Migration{
	ID:     "20260715021827",
	Before: "20260714033059",
	Upgrade: func(db *gorm.DB) error {
		db.Exec(
			"CREATE TYPE priority AS ENUM ('low', 'medium', 'high')",
		)
		db.Exec(
			"CREATE TABLE todos (" +
				"id SERIAL PRIMARY KEY," +
				"uuid UUID DEFAULT gen_random_uuid()," +
				"user_id INT NOT NULL," +
				"label_id INT DEFAULT NULL," +
				"title VARCHAR(255) NOT NULL," +
				"description TEXT," +
				"due_date DATE," +
				"start_date DATE," +
				"parent_id INT DEFAULT NULL," +
				"priority priority DEFAULT 'medium' NOT NULL," +
				"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
				"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
				"CONSTRAINT fk_todos_parent FOREIGN KEY (parent_id) REFERENCES todos(id) ON DELETE SET NULL" +
				")",
		)

		return nil
	},
	Downgrade: func(db *gorm.DB) error {
		db.Exec("DROP TABLE IF EXISTS todos")
		db.Exec("DROP TYPE IF EXISTS priority")
		return nil
	},
}

func init() {
	Register(Migration20260715021827)
}
