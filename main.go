package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/inner-daydream/kvz/kv"
	"github.com/inner-daydream/kvz/repositories/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlite.CreateDB("kv.db")
	if err != nil {
		log.Fatal(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}
	var repo kv.KvRepository = sqlite.New(db)
	var service kv.KvService = kv.NewServcice(repo)
	service.Set("path", "this/is/my/path")
	val, err := service.Get("path")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched val is: %s", val)
}
