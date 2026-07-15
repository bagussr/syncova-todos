package migrations

import "gorm.io/gorm"

var Migration20260715024000 = Migration{
	ID:     "20260715024000",
	Before: "20260715021827",
	Upgrade: func(db *gorm.DB) error {
		if err := db.Exec(
			"CREATE TABLE todos_statuses (" +
				"id SERIAL PRIMARY KEY," +
				"uuid UUID DEFAULT gen_random_uuid()," +
				"project_id INT NOT NULL," +
				"status VARCHAR(50) NOT NULL," +
				"CONSTRAINT fk_todos_statuses_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE" +
				")",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"CREATE TABLE labels_statuses (" +
				"id SERIAL PRIMARY KEY," +
				"uuid UUID DEFAULT gen_random_uuid()," +
				"project_id INT NOT NULL," +
				"label VARCHAR(50) NOT NULL," +
				"CONSTRAINT fk_labels_statuses_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE" +
				")",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos ADD CONSTRAINT fk_todos_label FOREIGN KEY (label_id) REFERENCES labels_statuses(id) ON DELETE SET NULL",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"CREATE TABLE todos_status_todos (" +
				"id SERIAL PRIMARY KEY," +
				"uuid UUID DEFAULT gen_random_uuid()," +
				"todos_status_id INT NOT NULL," +
				"todos_id INT NOT NULL," +
				"CONSTRAINT fk_todos_status_todos_status FOREIGN KEY (todos_status_id) REFERENCES todos_statuses(id) ON DELETE CASCADE," +
				"CONSTRAINT fk_todos_status_todos_todos FOREIGN KEY (todos_id) REFERENCES todos(id) ON DELETE CASCADE" +
				")",
		).Error; err != nil {
			return err
		}

		return nil
	},
	Downgrade: func(db *gorm.DB) error {
		if err := db.Exec("DROP TABLE IF EXISTS todos_status_todos").Error; err != nil {
			return err
		}
		if err := db.Exec("ALTER TABLE todos DROP CONSTRAINT IF EXISTS fk_todos_label").Error; err != nil {
			return err
		}
		if err := db.Exec("DROP TABLE IF EXISTS labels_statuses").Error; err != nil {
			return err
		}
		if err := db.Exec("DROP TABLE IF EXISTS todos_statuses").Error; err != nil {
			return err
		}
		return nil
	},
}

func init() {
	Register(Migration20260715024000)
}
