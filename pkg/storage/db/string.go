package db

import (
	"fmt"

	"github.com/dolfly/mydis/pkg/storage"
	"github.com/tidwall/redcon"
)

func (ms dbStorage) cmd_append(conn redcon.Conn, args [][]byte) {
	if len(args) != 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	kv := storage.Set{
		RKey: string(args[1]),
	}
	has, err := ms.eg.Get(&kv)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	var affected int64
	if !has {
		kv.Value = args[2]
		affected, err = ms.eg.Insert(kv)
	} else {
		kv.Value = append(kv.Value, args[2]...)
		affected, err = ms.eg.ID(kv.Id).Cols("value").Update(kv)
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if affected > 0 {
		conn.WriteInt(len(kv.Value))
	} else {
		conn.WriteNull()
	}
}

func (ms dbStorage) cmd_del(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_dump(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_exists(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_expire(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_expireat(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_keys(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_bitcount(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}

func (ms dbStorage) cmd_bitfield(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_bitop(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_bitpos(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_decr(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_decrby(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_getrange(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_getbit(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_getset(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_incr(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_incrby(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_incrbyfloat(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_mget(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_mset(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_msetnx(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_psetex(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_set(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn < 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	kv := storage.Set{
		RKey: string(args[1]),
	}
	has, err := ms.eg.Get(&kv)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	var affected int64
	if !has {
		kv.Value = args[1]
		affected, err = ms.eg.InsertOne(kv)
	} else {
		affected, err = ms.eg.ID(kv.Id).Cols("value").Update(&kv)
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if affected > 0 {
		conn.WriteString("OK")
	} else {
		conn.WriteNull()
	}
}
func (ms dbStorage) cmd_setbit(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_setex(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_setnx(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_setrange(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_strlen(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	kv := storage.Set{
		RKey: string(args[1]),
	}
	has, err := ms.eg.Get(&kv)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if has {
		conn.WriteInt(len(kv.Value))
	} else {
		conn.WriteInt(0)
	}
}
