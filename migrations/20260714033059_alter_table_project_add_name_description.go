package migrations

import "gorm.io/gorm"

var Migration20260714033059 = Migration{
	ID:     "20260714033059",
	Before: "20260714031753",
	Upgrade: func(db *gorm.DB) error {
		db.Exec(
			"ALTER TABLE projects add column name VARCHAR(255) NOT NULL",
		)
		db.Exec(
			"ALTER TABLE projects add column description TEXT",
		)
		return nil
	},
	Downgrade: func(db *gorm.DB) error {
		db.Exec(
			"ALTER TABLE projects drop column name",
		)
		db.Exec(
			"ALTER TABLE projects drop column description",
		)
		return nil
	},
}

func init() {
	Register(Migration20260714033059)
}
