package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type migration struct {
	name string
	sql  string
}

func (db *DB) Migrate(ctx context.Context) error {
	if err := db.createMigrationTable(ctx); err != nil {
		return err
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		applied, err := db.isMigrationApplied(ctx, migration.name)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		if err := db.runMigration(ctx, migration); err != nil {
			return err
		}

		fmt.Printf("applied migration: %s\n", migration.name)
	}

	return nil
}

func (db *DB) createMigrationTable(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`

	_, err := db.conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	return nil
}

func (db *DB) isMigrationApplied(ctx context.Context, name string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM schema_migrations
			WHERE name = $1
		)
	`

	var applied bool

	if err := db.conn.QueryRow(ctx, query, name).Scan(&applied); err != nil {
		return false, fmt.Errorf(
			"checking migration %s: %w",
			name,
			err,
		)
	}

	return applied, nil
}

func (db *DB) runMigration(ctx context.Context, migration migration) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"starting transaction for migration %s: %w",
			migration.name,
			err,
		)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, migration.sql); err != nil {
		return fmt.Errorf(
			"running migration %s: %w",
			migration.name,
			err,
		)
	}

	const recordMigration = `
		INSERT INTO schema_migrations (name)
		VALUES ($1)
	`

	if _, err := tx.Exec(ctx, recordMigration, migration.name); err != nil {
		return fmt.Errorf(
			"recording migration %s: %w",
			migration.name,
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"committing migration %s: %w",
			migration.name,
			err,
		)
	}

	return nil
}

func loadMigrations() ([]migration, error) {
	var migrations []migration

	err := fs.WalkDir(migrationFiles, "migrations", func(
		path string,
		entry fs.DirEntry,
		err error,
	) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".sql" {
			return nil
		}

		contents, err := migrationFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", path, err)
		}

		name := strings.TrimPrefix(path, "migrations/")

		migrations = append(migrations, migration{
			name: name,
			sql:  string(contents),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("loading migrations: %w", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].name < migrations[j].name
	})

	return migrations, nil
}