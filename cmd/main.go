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
	db, err := sqlite.OpenDB("kv.db")
	if err != nil {
		log.Fatal(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}
	var repo kv.KvRepository = sqlite.New(db)
	var service kv.KvService = kv.NewServcice(repo)
	cliCommands := cli.NewCli(&service)
	cli.ParseAndExecute(cliCommands)
	// service.Set("path", "this/is/my/path")
	// val, err := service.Get("path")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Fetched val is: %s", val)
}
