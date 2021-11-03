package db

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/dolfly/mydis/pkg/storage"
	"github.com/tidwall/redcon"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"
)

var dict = map[string]schemas.DBType{
	"mysql":    schemas.MYSQL,
	"sqlite":   schemas.SQLITE,
	"sqlite3":  schemas.SQLITE,
	"oracle":   schemas.ORACLE,
	"mssql":    schemas.MSSQL,
	"postgres": schemas.POSTGRES,
}

type dbStorage struct {
	eg *xorm.EngineGroup
	db string
}

func New(driver string, sources ...string) (s storage.Storage, err error) {
	ms := dbStorage{}
	ms.db = "0"
	ms.eg, err = xorm.NewEngineGroup(driver, sources, xorm.RandomPolicy())
	if err != nil {
		return
	}
	go func() {
		if len(ms.eg.Slaves()) > 0 {
			return
		}
		lock := sync.Mutex{}
		master := ms.eg.Master()
		ticker := time.NewTicker(5 * time.Second)
		for _ = range ticker.C {
			if err := master.Ping(); err != nil {
				eg, _ := xorm.NewEngineGroup(ms.eg.Slave(), ms.eg.Slaves())
				lock.Lock()
				ms.eg = eg
				lock.Unlock()
			} else {
				eg, _ := xorm.NewEngineGroup(master, ms.eg.Slaves())
				lock.Lock()
				ms.eg = eg
				lock.Unlock()
			}
		}
	}()

	ms.eg.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 100000000000))
	ms.eg.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	ms.eg.SetTableMapper(names.NewPrefixMapper(names.SameMapper{}, "t_"))
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

func (ms dbStorage) Debug() {
	ms.eg.SetLogLevel(log.LOG_DEBUG)
	ms.eg.SetLogger(log.NewSimpleLogger(os.Stdout))
	ms.eg.Logger().ShowSQL(true)
}

func (ms dbStorage) Handler(conn redcon.Conn, cmd redcon.Command) {
	sc := strings.ToLower(string(cmd.Args[0]))
	cmdfunc, ok := ms.command(sc)
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

func (ms dbStorage) command(name string) (func(conn redcon.Conn, args [][]byte), bool) {
	var commands = map[string]func(conn redcon.Conn, args [][]byte){
		// @system
		"ping":   ms.cmd_ping,
		"quit":   ms.cmd_quit,
		"dbtype": ms.cmd_dbtype,
		"dbdump": ms.cmd_dbdump,
		// @generic
		"del":      ms.cmd_del,      //TODO:
		"dump":     ms.cmd_dump,     //TODO:
		"exists":   ms.cmd_exists,   //TODO:
		"expire":   ms.cmd_expire,   //TODO:
		"expireat": ms.cmd_expireat, //TODO:
		"keys":     ms.cmd_keys,     //TODO:
		// @string
		"append":      ms.cmd_append,      //
		"bitcount":    ms.cmd_bitcount,    //TODO:
		"bitfield":    ms.cmd_bitfield,    //TODO:
		"bitop":       ms.cmd_bitop,       //TODO:
		"bitpos":      ms.cmd_bitpos,      //TODO:
		"decr":        ms.cmd_decr,        //
		"decrby":      ms.cmd_decrby,      //
		"get":         ms.cmd_get,         //TODO:
		"getbit":      ms.cmd_getbit,      //TODO:
		"getrange":    ms.cmd_getrange,    //TODO:
		"getset":      ms.cmd_getset,      //TODO:
		"incr":        ms.cmd_incr,        //
		"incrby":      ms.cmd_incrby,      //
		"incrbyfloat": ms.cmd_incrbyfloat, //
		"mget":        ms.cmd_mget,        //
		"mset":        ms.cmd_mset,        //
		"msetnx":      ms.cmd_msetnx,      //
		"psetex":      ms.cmd_psetex,      //TODO:
		"set":         ms.cmd_set,         //
		"setbit":      ms.cmd_setbit,      //TODO:
		"setex":       ms.cmd_setex,       //TODO:
		"setnx":       ms.cmd_setnx,       //TODO:
		"setrange":    ms.cmd_setrange,    //TODO:
		"strlen":      ms.cmd_strlen,      //TODO:
		// @Hash
		"hdel":         ms.cmd_hdel,
		"hexists":      ms.cmd_hexists,
		"hget":         ms.cmd_hget,
		"hgetall":      ms.cmd_hgetall,
		"hincrby":      ms.cmd_hincrby,
		"hincrbyfloat": ms.cmd_hincrbyfloat,
		"hkeys":        ms.cmd_hkeys,
		"hlen":         ms.cmd_hlen,
		"hmget":        ms.cmd_hmget,
		"hmset":        ms.cmd_hmset,
		"hscan":        ms.cmd_hscan,   //@TODO:
		"hset":         ms.cmd_hset,    //
		"hsetnx":       ms.cmd_hsetnx,  //@TODO:
		"hstrlen":      ms.cmd_hstrlen, //
		"hvals":        ms.cmd_hvals,   //@TODO:
		// @list
		"blpop":      ms.cmd_blpop,      // TODO:
		"brpop":      ms.cmd_brpop,      // TODO:
		"brpoplpush": ms.cmd_brpoplpush, // TODO:
		"lindex":     ms.cmd_lindex,     // TODO:
		"linsert":    ms.cmd_linsert,    // TODO:
		"llen":       ms.cmd_llen,       // TODO:
		"lpop":       ms.cmd_lpop,       // TODO:
		"lpush":      ms.cmd_lpush,      // TODO:
		"lpushx":     ms.cmd_lpushx,     // TODO:
		"lrange":     ms.cmd_lrange,     // TODO:
		"lrem":       ms.cmd_lrem,       // TODO:
		"lset":       ms.cmd_lset,       // TODO:
		"ltrim":      ms.cmd_ltrim,      // TODO:
		"rpop":       ms.cmd_rpop,       // TODO:
		"rpoplpush":  ms.cmd_rpoplpush,  // TODO:
		"rpush":      ms.cmd_rpush,      // TODO:
		"rpushx":     ms.cmd_rpushx,     // TODO:
	}
	comfunc, ok := commands[name]
	return comfunc, ok
}
