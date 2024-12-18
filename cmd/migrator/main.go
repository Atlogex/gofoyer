package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"path/filepath"
)

func main() {
	const databaseName = "gofoyer.db"
	var storagePath, migrationsPath, migrationsTable string

	dir, _ := filepath.Abs("")
	fmt.Println(dir)
	defaultDbPath := filepath.Join(dir, "database", databaseName)
	defaultMigrationPath := filepath.Join(dir, "migrations")

	flag.StringVar(&storagePath, "storage_path", defaultDbPath, "path to storage file")
	flag.StringVar(&migrationsPath, "migrations_path", defaultMigrationPath, "path to migration file")
	flag.StringVar(&migrationsTable, "migrations_table", "migrations", "name of migrations table")

	flag.Parse()

	fmt.Println("storagePath", storagePath)
	fmt.Println("migrationsPath", migrationsPath)

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes in database")
		} else {
			panic(err)
		}
	}

	fmt.Println("Database migrated")
}
