package migrations

import "gorm.io/gorm"

var Migration20260715072135 = Migration{
	ID:     "20260715072135",
	Before: "20260715024000",
	Upgrade: func(db *gorm.DB) error {
		if err := db.Exec(
			"ALTER TABLE projects ADD COLUMN IF NOT EXISTS user_id VARCHAR(255)",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"UPDATE projects SET user_id = '' WHERE user_id IS NULL",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE projects ALTER COLUMN user_id SET NOT NULL",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE projects ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE projects ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos ALTER COLUMN user_id TYPE VARCHAR(255) USING user_id::VARCHAR(255)",
		).Error; err != nil {
			return err
		}
		return nil
	},
	Downgrade: func(db *gorm.DB) error {
		if err := db.Exec(
			"ALTER TABLE projects DROP COLUMN IF EXISTS user_id, DROP COLUMN IF EXISTS created_at, DROP COLUMN IF EXISTS updated_at",
		).Error; err != nil {
			return err
		}

		if err := db.Exec(
			"ALTER TABLE todos ALTER COLUMN user_id TYPE INT USING user_id::INT",
		).Error; err != nil {
			return err
		}
		return nil
	},
}

func init() {
	Register(Migration20260715072135)
}
