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
	eg       *xorm.EngineGroup
	db       string
	commands map[string]func(conn redcon.Conn, args [][]byte)
}

func (ms dbStorage) Handler(conn redcon.Conn, cmd redcon.Command) {
	sc := strings.ToLower(string(cmd.Args[0]))
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
func (ms dbStorage) cmd(name string) func(conn redcon.Conn, args [][]byte) {
	return func(conn redcon.Conn, args [][]byte) {
		conn.WriteString("OK")
	}
}
func (ms dbStorage) cmd_ping(conn redcon.Conn, args [][]byte) {
	conn.WriteString("PONG")
}
func (ms dbStorage) cmd_quit(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
	conn.Close()
}
func (ms dbStorage) cmd_get(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	set := storage.Set{RKey: string(args[1])}
	has, err := ms.eg.Get(&set)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if has {
		conn.WriteBulk([]byte(set.Value))
	} else {
		conn.WriteNull()
	}

}
func (ms dbStorage) cmd_db(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func Commands(ms dbStorage) map[string]func(conn redcon.Conn, args [][]byte) {
	return map[string]func(conn redcon.Conn, args [][]byte){
		"ping": ms.cmd_ping,
		"quit": ms.cmd_quit,
		"db":   ms.cmd_db,
		// @generic
		"del":      ms.cmd_del,
		"dump":     ms.cmd_dump,
		"exists":   ms.cmd_exists,
		"expire":   ms.cmd_expire,
		"expireat": ms.cmd_expireat,
		"keys":     ms.cmd_keys,
		// @string
		"append":      ms.cmd_append,
		"bitcount":    ms.cmd_bitcount,
		"bitfield":    ms.cmd_bitfield,
		"bitop":       ms.cmd_bitop,
		"bitpos":      ms.cmd_bitpos,
		"decr":        ms.cmd_decr,
		"decrby":      ms.cmd_decrby,
		"get":         ms.cmd_get,
		"getbit":      ms.cmd_getbit,
		"getrange":    ms.cmd_getrange,
		"getset":      ms.cmd_getset,
		"incr":        ms.cmd_incr,
		"incrby":      ms.cmd_incrby,
		"incrbyfloat": ms.cmd_incrbyfloat,
		"mget":        ms.cmd_mget,
		"mset":        ms.cmd_mset,
		"msetnx":      ms.cmd_msetnx,
		"psetex":      ms.cmd_psetex,
		"set":         ms.cmd_set,
		"setbit":      ms.cmd_setbit,
		"setex":       ms.cmd_setex,
		"setnx":       ms.cmd_setnx,
		"setrange":    ms.cmd_setrange,
		"strlen":      ms.cmd_strlen,
		// @Hash
		"hdel":         ms.cmd_hdel,
		"hexists":      ms.cmd_hexists,
		"hget":         ms.cmd_hget,
		"hgetall":      ms.cmd_hgetall,
		"hincrby":      ms.cmd_hincrby,
		"hincrbyfloat": ms.cmd_hincrbyfloat,
		"hkeys":        ms.cmd_hkeys,
		"hlen":         ms.cmd_hlen,    //@TODO:
		"hmget":        ms.cmd_hmget,   //@TODO:
		"hmset":        ms.cmd_hmset,   //@TODO:
		"hscan":        ms.cmd_hscan,   //@TODO:
		"hset":         ms.cmd_hset,    //
		"hsetnx":       ms.cmd_hsetnx,  //@TODO:
		"hsetlen":      ms.cmd_hsetlen, //@TODO:
		"hvals":        ms.cmd_hvals,   //@TODO:
		// @list
	}
}
func New(driver string, sources ...string) (s storage.Storage, err error) {
	ms := dbStorage{}
	ms.db = "0"
	ms.eg, err = xorm.NewEngineGroup(driver, sources, xorm.RandomPolicy())
	if err != nil {
		return
	}
	ms.eg.SetLogLevel(log.LOG_DEBUG)
	ms.eg.SetLogger(log.NewSimpleLogger(os.Stdout))
	ms.eg.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 1000))
	ms.eg.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	ms.eg.SetTableMapper(names.NewPrefixMapper(names.SameMapper{}, "t_"))
	ms.commands = Commands(ms)
	err = ms.eg.Sync2(
		new(storage.Namespace),
		new(storage.Set),
		new(storage.List),
		new(storage.Hash),
		new(storage.ZSet),
	)
	ms.eg.Insert(&storage.Namespace{
		Name: ms.db,
		RKey: "mydis",
		Type: 1,
	})
	if err != nil {
		return
	}
	return &ms, err
}
