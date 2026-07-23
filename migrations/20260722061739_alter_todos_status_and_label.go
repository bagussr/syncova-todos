package migrations

import "gorm.io/gorm"

var Migration20260722061739 = Migration{
	ID:     "20260722061739",
	Before: "20260715072135",
	Upgrade: func(db *gorm.DB) error {
		if err := db.Exec(
			"ALTER TABLE todos ADD COLUMN IF NOT EXISTS project_id INT NULL",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_statuses ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_statuses ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE labels_statuses ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE labels_statuses ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_status_todos ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_status_todos ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos ADD CONSTRAINT fk_todos_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL",
		).Error; err != nil {
			return err
		}

		return nil
	},
	Downgrade: func(db *gorm.DB) error {
		if err := db.Exec(
			"ALTER TABLE todos DROP CONSTRAINT IF EXISTS fk_todos_project",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos DROP COLUMN IF EXISTS project_id",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_statuses DROP COLUMN IF EXISTS created_at",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_statuses DROP COLUMN IF EXISTS updated_at",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE labels_statuses DROP COLUMN IF EXISTS created_at",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE labels_statuses DROP COLUMN IF EXISTS updated_at",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_status_todos DROP COLUMN IF EXISTS created_at",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos_status_todos DROP COLUMN IF EXISTS updated_at",
		).Error; err != nil {
			return err
		}

		return nil
	},
}

func init() {
	Register(Migration20260722061739)
}
