package main

import (
	"log"

	"github.com/dolfly/mydis/pkg/storage"
	"github.com/dolfly/mydis/pkg/storage/db"
	"github.com/tidwall/redcon"
)

var addr = ":6380"

var s storage.Storage

func init() {
	var err error
	s, err = db.New("sqlite3", "mydis.db", "backup.db")
	if err != nil {
		panic(err)
	}
}

func main() {
	go log.Printf("started server at %s", addr)
	err := redcon.ListenAndServe(addr, s.Handler, s.Accept, s.Closed)
	if err != nil {
		log.Fatal(err)
	}
}
