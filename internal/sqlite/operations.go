package sqlite

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %w", err)
	}
	return db, nil
}

// update the db to the latest sql schema
func Migrate(db *sql.DB) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetDialect("sqlite3")
	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(db, "sql/migrations"); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}
	return nil
}

// func isOutdated(db *sql.DB, lastKnownVersion int64) (bool, error) {
// 	currentVersion, err := goose.GetDBVersion(db)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to fetch db version: %w", err)
// 	}
// 	isOutdated := lastKnownVersion < currentVersion
// 	return isOutdated, nil
// }
