package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syncova-todo/config"
	"syncova-todo/infrastructure/database"
	"syncova-todo/migrations"
	"time"

	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := strings.ToLower(strings.TrimSpace(os.Args[1]))
	switch command {
	case "up":
		runUp()
	case "down":
		runDown(os.Args[2:])
	case "create":
		runCreate(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func runUp() {
	cfg := config.LoadConfig()

	postgres, err := database.NewPostgresConnection(cfg, false)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		if err := postgres.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	if err := migrations.BackfillHistoricalBeforeIDs(postgres.DB); err != nil {
		log.Fatalf("failed to backfill historical migrations: %v", err)
	}

	if err := migrations.ApplyAll(postgres.DB); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
}

func runDown(args []string) {
	cfg := config.LoadConfig()

	postgres, err := database.NewPostgresConnection(cfg, false)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		if err := postgres.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	if err := executeDown(postgres.DB, args); err != nil {
		log.Fatalf("rollback failed: %v", err)
	}

	log.Println("migration rollback completed successfully")
}

func executeDown(db *gorm.DB, args []string) error {
	if len(args) == 0 {
		return migrations.RollbackLast(db)
	}

	first := strings.TrimSpace(args[0])
	if first == "--steps" || first == "-s" {
		if len(args) < 2 {
			return fmt.Errorf("missing value for %s", first)
		}
		steps, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			return fmt.Errorf("invalid steps value %q", args[1])
		}
		return migrations.RollbackSteps(db, steps)
	}

	return migrations.RollbackTo(db, first)
}

func runCreate(args []string) {
	if len(args) < 1 {
		fmt.Println("missing migration name")
		fmt.Println("usage: go run ./cmd/migrations create <name> [before_migration_id]")
		os.Exit(1)
	}

	name := sanitizeName(args[0])
	if name == "" {
		log.Fatal("invalid migration name")
	}

	beforeID := ""
	if len(args) >= 2 {
		beforeID = strings.TrimSpace(args[1])
	}
	if beforeID == "" {
		beforeID = detectBeforeMigrationID()
	}

	migrationID := time.Now().UTC().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.go", migrationID, name)
	filePath := filepath.Join("migrations", fileName)

	if _, err := os.Stat(filePath); err == nil {
		log.Fatalf("migration file already exists: %s", filePath)
	}

	content := migrationTemplate(migrationID, beforeID)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Fatalf("failed to create migration file: %v", err)
	}

	fmt.Printf("created migration: %s\n", filePath)
}

func detectBeforeMigrationID() string {
	if id := latestAppliedMigrationFromDB(); id != "" {
		return id
	}

	return latestMigrationIDFromFiles()
}

func latestAppliedMigrationFromDB() string {
	cfg := config.LoadConfig()
	postgres, err := database.NewPostgresConnection(cfg, false)
	if err != nil {
		return ""
	}
	defer func() {
		_ = postgres.Close()
	}()

	tableName := ""
	type tableResult struct {
		Name string
	}
	var tbl tableResult
	if err := postgres.DB.Raw(`
		SELECT table_name AS name
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name IN ('migrations_table', 'migrations')
		ORDER BY CASE WHEN table_name = 'migrations_table' THEN 0 ELSE 1 END
		LIMIT 1
	`).Scan(&tbl).Error; err != nil {
		return ""
	}
	tableName = strings.TrimSpace(tbl.Name)
	if tableName == "" {
		return ""
	}

	type latest struct {
		MigrationID string
	}
	var result latest
	query := fmt.Sprintf("SELECT migration_id FROM public.%s ORDER BY applied_at DESC, id DESC LIMIT 1", tableName)
	if err := postgres.DB.Raw(query).Scan(&result).Error; err != nil {
		return ""
	}

	return strings.TrimSpace(result.MigrationID)
}

func latestMigrationIDFromFiles() string {
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return ""
	}

	idPrefix := regexp.MustCompile(`^(\d{14})_`)
	ids := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || name == "main.go" {
			continue
		}
		m := idPrefix.FindStringSubmatch(name)
		if len(m) < 2 {
			continue
		}
		ids = append(ids, m[1])
	}

	if len(ids) == 0 {
		return ""
	}

	sort.Strings(ids)
	return ids[len(ids)-1]
}

func migrationTemplate(migrationID, beforeID string) string {
	varName := "Migration" + migrationID
	return fmt.Sprintf(`package migrations

import "gorm.io/gorm"

var %s = Migration{
	ID:     %q,
	Before: %q,
	Upgrade: func(db *gorm.DB) error {
		return nil
	},
	Downgrade: func(db *gorm.DB) error {
		return nil
	},
}

func init() {
	Register(%s)
}
`, varName, migrationID, beforeID, varName)
}

func sanitizeName(name string) string {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	cleaned = strings.ReplaceAll(cleaned, " ", "_")
	re := regexp.MustCompile(`[^a-z0-9_]+`)
	cleaned = re.ReplaceAllString(cleaned, "")
	cleaned = strings.Trim(cleaned, "_")
	return cleaned
}

func printUsage() {
	fmt.Println("migration command")
	fmt.Println("usage:")
	fmt.Println("  go run ./cmd/migrations up")
	fmt.Println("  go run ./cmd/migrations down")
	fmt.Println("  go run ./cmd/migrations down --steps <n>")
	fmt.Println("  go run ./cmd/migrations down <migration_id>")
	fmt.Println("  go run ./cmd/migrations create <name> [before_migration_id]")
}
