package main

import (
	_ "embed"
	"log"

	"github.com/inner-daydream/kvz/internal/cli"
	"github.com/inner-daydream/kvz/internal/kv"
	"github.com/inner-daydream/kvz/internal/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: store the db in a better/user defined place
	dbPath := "kv.db"
	db, err := sqlite.OpenDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlite.Migrate(db)
	queries := sqlite.New(db)
	repo := sqlite.NewRepository(queries)
	var service kv.KvService = kv.NewServcice(repo)
	if err != nil {
		log.Fatal(err)
	}
	cliCommands := cli.NewCli(service)
	cli.ParseAndExecute(cliCommands)
}
