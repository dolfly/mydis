package mysql

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/dolfly/mydis/pkg/store"
	"github.com/tidwall/redcon"
)

type mysqlStore struct {
	db       *sqlx.DB
	commands map[string]func(conn redcon.Conn, args [][]byte)
}

func (ms mysqlStore) Handler(conn redcon.Conn, cmd redcon.Command) {
	sc := strings.ToLower(string(cmd.Args[0]))
	cmdfunc, ok := ms.commands[sc]
	if ok {
		cmdfunc(conn, cmd.Args)
	}
}
func (ms mysqlStore) Accept(conn redcon.Conn) bool {
	return true
}
func (ms mysqlStore) Closed(conn redcon.Conn, err error) {
}
func (ms mysqlStore) get(conn redcon.Conn, args [][]byte) {
}
func (ms mysqlStore) quit(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
	conn.Close()
}
func (ms mysqlStore) ping(conn redcon.Conn, args [][]byte) {
	conn.WriteString("PONG")
}
func (ms mysqlStore) Commands() map[string]func(conn redcon.Conn, args [][]byte) {
	return map[string]func(conn redcon.Conn, args [][]byte){
		"get":  ms.get,
		"ping": ms.ping,
	}
}

func New(dsn string) (s store.Store, err error) {
	ms := mysqlStore{}
	ms.commands = ms.Commands()
	ms.db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	return &ms, err
}
