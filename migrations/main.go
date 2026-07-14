package migrations

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	ID        string
	Before    string
	Upgrade   func(db *gorm.DB) error
	Downgrade func(db *gorm.DB) error
}

type MigrationRecord struct {
	ID                uint      `gorm:"primaryKey;autoIncrement"`
	MigrationID       string    `gorm:"column:migration_id;type:varchar(255);not null;uniqueIndex"`
	BeforeMigrationID string    `gorm:"column:before_migration_id;type:varchar(255)"`
	AppliedAt         time.Time `gorm:"column:applied_at;not null;autoCreateTime"`
}

func (MigrationRecord) TableName() string {
	return "migrations"
}

var registry []Migration

func Register(m Migration) {
	if m.ID == "" {
		panic("migration id cannot be empty")
	}
	if m.Upgrade == nil {
		panic(fmt.Sprintf("migration %q upgrade function cannot be nil", m.ID))
	}

	for _, existing := range registry {
		if existing.ID == m.ID {
			panic(fmt.Sprintf("migration id %q is already registered", m.ID))
		}
	}

	registry = append(registry, m)
}

func Registered() []Migration {
	items := make([]Migration, len(registry))
	copy(items, registry)
	return items
}

func ApplyAll(db *gorm.DB) error {
	if db == nil {
		return errors.New("db cannot be nil")
	}

	if err := ensureMigrationTable(db); err != nil {
		return err
	}

	if err := BackfillHistoricalBeforeIDs(db); err != nil {
		return err
	}

	ordered, err := orderedMigrations(registry)
	if err != nil {
		return err
	}

	appliedMap, err := loadAppliedMap(db)
	if err != nil {
		return err
	}

	for _, migration := range ordered {
		if appliedMap[migration.ID] {
			continue
		}

		if migration.Before != "" && !appliedMap[migration.Before] {
			return fmt.Errorf("migration %q depends on previous migration %q, but %q is not applied", migration.ID, migration.Before, migration.Before)
		}

		beforeIDForRecord := migration.Before
		if beforeIDForRecord == "" {
			lastAppliedID, err := getLastAppliedMigrationID(db)
			if err != nil {
				return err
			}
			beforeIDForRecord = lastAppliedID
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Upgrade(tx); err != nil {
				return fmt.Errorf("upgrade %q failed: %w", migration.ID, err)
			}

			record := MigrationRecord{
				MigrationID:       migration.ID,
				BeforeMigrationID: beforeIDForRecord,
			}

			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("failed to persist migration %q: %w", migration.ID, err)
			}

			return nil
		}); err != nil {
			return err
		}

		appliedMap[migration.ID] = true
	}

	return nil
}

func RollbackLast(db *gorm.DB) error {
	return RollbackSteps(db, 1)
}

func RollbackSteps(db *gorm.DB, steps int) error {
	if db == nil {
		return errors.New("db cannot be nil")
	}
	if steps <= 0 {
		return errors.New("steps must be greater than 0")
	}

	if err := ensureMigrationTable(db); err != nil {
		return err
	}

	tableName, err := resolveMigrationTableName(db)
	if err != nil {
		return err
	}
	if tableName == "" {
		return errors.New("no migrations table found in public schema")
	}

	records, err := loadAppliedMigrationsDesc(db, tableName)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errors.New("no applied migration found to rollback")
	}
	if steps > len(records) {
		return fmt.Errorf("cannot rollback %d step(s), only %d migration(s) applied", steps, len(records))
	}

	lookup := buildMigrationLookup()
	for i := 0; i < steps; i++ {
		if err := rollbackRecord(db, tableName, records[i], lookup); err != nil {
			return err
		}
	}

	return nil
}

func RollbackTo(db *gorm.DB, migrationID string) error {
	if db == nil {
		return errors.New("db cannot be nil")
	}
	targetID := strings.TrimSpace(migrationID)
	if targetID == "" {
		return errors.New("migration id cannot be empty")
	}

	if err := ensureMigrationTable(db); err != nil {
		return err
	}

	tableName, err := resolveMigrationTableName(db)
	if err != nil {
		return err
	}
	if tableName == "" {
		return errors.New("no migrations table found in public schema")
	}

	records, err := loadAppliedMigrationsDesc(db, tableName)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errors.New("no applied migration found to rollback")
	}

	targetIndex := -1
	for idx, record := range records {
		if strings.TrimSpace(record.MigrationID) == targetID {
			targetIndex = idx
			break
		}
	}
	if targetIndex == -1 {
		return fmt.Errorf("target migration %q is not applied", targetID)
	}
	if targetIndex == 0 {
		return nil
	}

	lookup := buildMigrationLookup()
	for i := 0; i < targetIndex; i++ {
		if err := rollbackRecord(db, tableName, records[i], lookup); err != nil {
			return err
		}
	}

	return nil
}

type appliedMigration struct {
	ID          uint
	MigrationID string
}

func loadAppliedMigrationsDesc(db *gorm.DB, tableName string) ([]appliedMigration, error) {
	query := fmt.Sprintf("SELECT id, migration_id FROM public.%s ORDER BY applied_at DESC, id DESC", tableName)
	var records []appliedMigration
	if err := db.Raw(query).Scan(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to load applied migrations: %w", err)
	}
	return records, nil
}

func buildMigrationLookup() map[string]Migration {
	lookup := make(map[string]Migration, len(registry))
	for _, m := range registry {
		lookup[m.ID] = m
	}
	return lookup
}

func rollbackRecord(db *gorm.DB, tableName string, record appliedMigration, lookup map[string]Migration) error {
	migrationID := strings.TrimSpace(record.MigrationID)
	if record.ID == 0 || migrationID == "" {
		return errors.New("invalid migration record to rollback")
	}

	migration, found := lookup[migrationID]
	if !found {
		return fmt.Errorf("migration %q is applied in database but not registered in code", migrationID)
	}
	if migration.Downgrade == nil {
		return fmt.Errorf("migration %q has no downgrade function", migrationID)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := migration.Downgrade(tx); err != nil {
			return fmt.Errorf("downgrade %q failed: %w", migrationID, err)
		}

		deleteQuery := fmt.Sprintf("DELETE FROM public.%s WHERE id = ?", tableName)
		if err := tx.Exec(deleteQuery, record.ID).Error; err != nil {
			return fmt.Errorf("failed to delete migration record %q: %w", migrationID, err)
		}

		return nil
	})
}

func BackfillHistoricalBeforeIDs(db *gorm.DB) error {
	if db == nil {
		return errors.New("db cannot be nil")
	}

	tableName, err := resolveMigrationTableName(db)
	if err != nil {
		return err
	}
	if tableName == "" {
		return nil
	}

	query := fmt.Sprintf("SELECT id, migration_id, before_migration_id FROM public.%s ORDER BY applied_at ASC, id ASC", tableName)
	type row struct {
		ID        uint
		Migration string
		Before    string
	}

	var rows []row
	if err := db.Raw(query).Scan(&rows).Error; err != nil {
		return fmt.Errorf("failed to read historical migrations from public.%s: %w", tableName, err)
	}

	for i := 1; i < len(rows); i++ {
		if strings.TrimSpace(rows[i].Before) != "" {
			continue
		}

		prevID := strings.TrimSpace(rows[i-1].Migration)
		if prevID == "" {
			continue
		}

		update := fmt.Sprintf("UPDATE public.%s SET before_migration_id = ? WHERE id = ?", tableName)
		if err := db.Exec(update, prevID, rows[i].ID).Error; err != nil {
			return fmt.Errorf("failed to backfill before_migration_id for row %d: %w", rows[i].ID, err)
		}
	}

	return nil
}

func resolveMigrationTableName(db *gorm.DB) (string, error) {
	type tbl struct {
		Name string
	}

	var result tbl
	if err := db.Raw(`
		SELECT table_name AS name
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name IN ('migrations_table', 'migrations')
		ORDER BY CASE WHEN table_name = 'migrations_table' THEN 0 ELSE 1 END
		LIMIT 1
	`).Scan(&result).Error; err != nil {
		return "", fmt.Errorf("failed to resolve migrations table in public schema: %w", err)
	}

	return result.Name, nil
}

func getLastAppliedMigrationID(db *gorm.DB) (string, error) {
	tableName, err := resolveMigrationTableName(db)
	if err != nil {
		return "", err
	}
	if tableName == "" {
		return "", nil
	}

	type result struct {
		MigrationID string
	}

	query := fmt.Sprintf("SELECT migration_id FROM public.%s ORDER BY applied_at DESC, id DESC LIMIT 1", tableName)
	var latest result
	if err := db.Raw(query).Scan(&latest).Error; err != nil {
		return "", fmt.Errorf("failed to read last applied migration from public.%s: %w", tableName, err)
	}

	return strings.TrimSpace(latest.MigrationID), nil
}

func ensureMigrationTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to migrate migrations table: %w", err)
	}
	return nil
}

func loadAppliedMap(db *gorm.DB) (map[string]bool, error) {
	var applied []MigrationRecord
	if err := db.Find(&applied).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	result := make(map[string]bool, len(applied))
	for _, item := range applied {
		result[item.MigrationID] = true
	}

	return result, nil
}

func orderedMigrations(items []Migration) ([]Migration, error) {
	if len(items) == 0 {
		return nil, nil
	}

	byID := make(map[string]Migration, len(items))
	adj := make(map[string][]string, len(items))
	inDegree := make(map[string]int, len(items))

	for _, item := range items {
		if _, exists := byID[item.ID]; exists {
			return nil, fmt.Errorf("duplicate migration id %q", item.ID)
		}
		byID[item.ID] = item
		inDegree[item.ID] = 0
	}

	for _, item := range items {
		if item.Before == "" {
			continue
		}
		if _, exists := byID[item.Before]; !exists {
			return nil, fmt.Errorf("migration %q references unknown previous migration id %q", item.ID, item.Before)
		}
		adj[item.Before] = append(adj[item.Before], item.ID)
		inDegree[item.ID]++
	}

	var queue []string
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue)

	ordered := make([]Migration, 0, len(items))
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		ordered = append(ordered, byID[id])

		next := adj[id]
		sort.Strings(next)
		for _, target := range next {
			inDegree[target]--
			if inDegree[target] == 0 {
				queue = append(queue, target)
				sort.Strings(queue)
			}
		}
	}

	if len(ordered) != len(items) {
		return nil, errors.New("cyclic migration dependencies detected")
	}

	return ordered, nil
}
