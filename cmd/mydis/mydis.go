package main

import (
	"flag"
	"log"

	"github.com/dolfly/mydis/pkg/storage/db"
	"github.com/tidwall/redcon"
)

var (
	address = flag.String("address", ":6380", "set server address")
	driver  = flag.String("driver", "sqlite3", "set db driver")
	source  = flag.String("source", "mydis.db", "set db source")
)

func main() {
	flag.Parse()
	s, err := db.New(*driver, *source)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = redcon.ListenAndServe(*address, s.Handler, s.Accept, s.Closed)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("started server at %s", *address)
}
