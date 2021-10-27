package db

import (
	"fmt"
	"os"
	"strings"
	"time"

	"xorm.io/xorm"
	"xorm.io/xorm/caches"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"

	"github.com/dolfly/mydis/pkg/storage"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/redcon"
)

type dbStorage struct {
	engine   *xorm.Engine
	db       string
	commands map[string]func(conn redcon.Conn, args [][]byte)
}

func (ms dbStorage) Handler(conn redcon.Conn, cmd redcon.Command) {
	sc := strings.ToLower(string(cmd.Args[0]))
	fmt.Println(cmd.Args)
	cmdfunc, ok := ms.commands[sc]
	if ok {
		cmdfunc(conn, cmd.Args)
	} else {
		conn.WriteError("ERR unknown command '" + sc + "'")
	}
}
func (ms dbStorage) Accept(conn redcon.Conn) bool {
	// Use this function to accept or deny the connection.
	fmt.Printf("accept: %s \n", conn.RemoteAddr())
	return true
}
func (ms dbStorage) Closed(conn redcon.Conn, err error) {
	// This is called when the connection has been closed
	fmt.Printf("closed: %s, err: %v \n", conn.RemoteAddr(), err)
}

func (ms dbStorage) cmd_quit(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
	conn.Close()
}
func (ms dbStorage) cmd_ping(conn redcon.Conn, args [][]byte) {
	conn.WriteString("PONG")
}
func (ms dbStorage) cmd_get(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	kv := storage.KV{RKey: string(args[1])}
	has, err := ms.engine.Get(&kv)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if has {
		conn.WriteBulk([]byte(kv.Value))
	} else {
		conn.WriteNull()
	}

}
func (ms dbStorage) cmd_set(conn redcon.Conn, args [][]byte) {
}
func (ms *dbStorage) cmd_db(conn redcon.Conn, args [][]byte) {
}
func (ms dbStorage) Commands() map[string]func(conn redcon.Conn, args [][]byte) {
	return map[string]func(conn redcon.Conn, args [][]byte){
		"ping": ms.cmd_ping,
		"quit": ms.cmd_quit,
		"get":  ms.cmd_get,
		"set":  ms.cmd_set,
		"db":   ms.cmd_db,
	}
}

func New(driver string, source string) (s storage.Storage, err error) {
	ms := dbStorage{}
	ms.commands = ms.Commands()
	ms.db = "0"
	ms.engine, err = xorm.NewEngine(driver, source)
	if err != nil {
		return
	}
	ms.engine.SetLogLevel(log.LOG_DEBUG)
	ms.engine.SetLogger(log.NewSimpleLogger(os.Stdout))
	ms.engine.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 1000))
	ms.engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	ms.engine.SetTableMapper(names.NewPrefixMapper(names.SameMapper{}, "t_"))
	err = ms.engine.Sync2(new(storage.Namespace), new(storage.KV), new(storage.Hash), new(storage.ZSet))
	ms.engine.Insert(&storage.Namespace{
		Name: ms.db,
		RKey: "mydis",
		Type: 1,
	})
	if err != nil {
		return
	}
	return &ms, err
}
