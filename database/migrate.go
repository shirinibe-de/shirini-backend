package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			migration_name TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create migration history table:", err)
	}

	entries, err := os.ReadDir("database/migrations")
	if err != nil {
		log.Fatal("Failed to read migrations:", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	for _, name := range migrationFiles {
		var alreadyApplied bool
		err := DB.QueryRow(
			`SELECT EXISTS (
				SELECT 1 FROM schema_migrations WHERE migration_name = $1
			)`,
			name,
		).Scan(&alreadyApplied)
		if err != nil {
			log.Fatal("Failed to check migration history:", err)
		}

		if alreadyApplied {
			log.Printf("Migration %s already applied; skipping", name)
			continue
		}

		path := filepath.Join("database/migrations", name)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read migration %s: %v", name, err)
		}

		tx, err := DB.Begin()
		if err != nil {
			log.Fatal("Failed to start migration transaction:", err)
		}

		if _, err = tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			log.Fatalf("Migration %s failed: %v", name, err)
		}

		_, err = tx.Exec(
			`INSERT INTO schema_migrations (migration_name) VALUES ($1)`,
			name,
		)
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("Failed to record migration %s: %v", name, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("Failed to commit migration %s: %v", name, err)
		}

		log.Printf("Migration %s completed successfully", name)
	}
}

func ValidateMigrationDirectory() error {
	entries, err := os.ReadDir("database/migrations")
	if err != nil {
		return fmt.Errorf("read migration directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		if strings.TrimSpace(entry.Name()) == "" {
			return fmt.Errorf("migration filename cannot be empty")
		}
	}

	return nil
}
