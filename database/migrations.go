package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migrate runs all pending migrations in ./migrations
func Migrate() {
	// Read migrations folder
	files, err := os.ReadDir("./migrations")
	if err != nil {
		log.Fatal("Could not read migrations folder:", err)
	}

	// Collect .sql files
	var migrationFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}

	// Sort alphabetically (from 000 to 006)
	sort.Strings(migrationFiles)

	// Apply pending migrations
	for _, file := range migrationFiles {
		var exists int
		err := DB.QueryRow("SELECT COUNT(1) FROM migration_history WHERE filename = ?", file).Scan(&exists)
		if err != nil {
			// Special case: if migration_history doesn’t exist yet,
			// this must be 000_create_migration_history.sql
			if strings.HasPrefix(file, "000_") {
				// Run it without checking history
				sqlBytes, readErr := os.ReadFile(filepath.Join("./migrations", file))
				if readErr != nil {
					log.Fatalf("Failed to read migration %s: %v", file, readErr)
				}
				sqlText := string(sqlBytes)

				_, execErr := DB.Exec(sqlText)
				if execErr != nil {
					log.Fatalf("Failed to execute migration %s: %v", file, execErr)
				}

				fmt.Println("✅ Applied bootstrap migration:", file)
				// Now insert it into the history table
				_, insertErr := DB.Exec("INSERT INTO migration_history (filename) VALUES (?)", file)
				if insertErr != nil {
					log.Fatalf("Failed to update migration history for %s: %v", file, insertErr)
				}
				continue
			}
			log.Fatal("Failed checking migration history:", err)
		}

		if exists > 0 {
			continue // already applied
		}

		// Read SQL file
		sqlBytes, err := os.ReadFile(filepath.Join("./migrations", file))
		if err != nil {
			log.Fatalf("Failed to read migration %s: %v", file, err)
		}
		sqlText := string(sqlBytes)

		// Execute migration
		_, err = DB.Exec(sqlText)
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}

		// Record in migration history
		_, err = DB.Exec("INSERT INTO migration_history (filename) VALUES (?)", file)
		if err != nil {
			log.Fatalf("Failed to update migration history for %s: %v", file, err)
		}

		fmt.Println("✅ Applied migration:", file)
	}
}
