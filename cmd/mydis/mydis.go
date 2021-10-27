package main

import (
	"log"

	"github.com/dolfly/mydis/pkg/store"
	"github.com/dolfly/mydis/pkg/store/mysql"
	"github.com/tidwall/redcon"
)

var addr = ":6380"

var s store.Store

func init(){
	s =  mysql.New()
}

func main() {
	go log.Printf("started server at %s", addr)
	err := redcon.ListenAndServe(addr, s.Handler, s.Accept, s.Closed)
	if err != nil {
		log.Fatal(err)
	}
}
