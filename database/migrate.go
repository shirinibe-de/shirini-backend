
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
	migrationsDir := "database/migrations"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatal("Failed to read migrations directory:", err)
	}

	var migrationFiles []string

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		migrationFiles = append(migrationFiles, entry.Name())
	}

	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		path := filepath.Join(migrationsDir, file)

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("Failed to read migration ", file, ":", err)
		}

		if _, err := DB.Exec(string(sqlBytes)); err != nil {
			log.Fatal("Migration failed for ", file, ":", err)
		}

		log.Println("Migration completed:", file)
	}

	if len(migrationFiles) == 0 {
		log.Fatal(fmt.Errorf("no migration files found in %s", migrationsDir))
	}
}
