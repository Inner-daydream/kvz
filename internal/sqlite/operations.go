package sqlite

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func OpenDB(path string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %w", err)
	}
	return db, nil
}

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

// update the db to the latest sql schema
func (m *SqliteMigrator) Migrate() error {
	goose.SetLogger(goose.NopLogger())
	goose.SetDialect("sqlite3")
	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(m.db, "sql/migrations"); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}
	return nil
}

type SqliteMigrator struct {
	db *sql.DB
}

func NewSqliteMigrator(db *sql.DB) *SqliteMigrator {
	return &SqliteMigrator{
		db: db,
	}
}
