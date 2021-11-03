package db

import (
	"fmt"
	"strconv"
	"strings"

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
	argn := len(args)
	if argn != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	args = append(args, []byte("1"))
	ms.cmd_decrby(conn, args)
}
func (ms dbStorage) cmd_decrby(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn != 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	valinc, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR incr value(%s) is not an integer or out of range", string(args[3])))
		return
	}
	kv := storage.Set{
		RKey: string(args[1]),
	}
	has, err := ms.eg.Cols("id", "rkey", "value").Get(&kv)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	var nval int64
	var value int64
	if has {
		value, err = strconv.ParseInt(string(kv.Value), 10, 64)
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR old value(%s) is not an integer or out of range", string(kv.Value)))
			return
		}
		nval = value - valinc
		kv.Value = []byte(fmt.Sprintf("%d", nval))
		_, err = ms.eg.ID(kv.Id).Update(&kv)
	} else {
		nval = -valinc
		kv.Value = []byte(fmt.Sprintf("%d", nval))
		_, err = ms.eg.Insert(&kv)
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteInt64(nval)
}
func (ms dbStorage) cmd_get(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	set := storage.Set{RKey: string(args[1])}
	has, err := ms.eg.Cols("id", "rkey", "value").Get(&set)
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
	argn := len(args)
	if argn != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	args = append(args, []byte("-1"))
	ms.cmd_decrby(conn, args)
}
func (ms dbStorage) cmd_incrby(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn != 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	valinc, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR incr value(%s) is not an integer or out of range", string(args[3])))
		return
	}
	args[2] = []byte(fmt.Sprintf("%d", 0-valinc))
	ms.cmd_decrby(conn, args)
}
func (ms dbStorage) cmd_incrbyfloat(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn != 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	valinc, err := strconv.ParseFloat(string(args[2]), 64)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR incr value(%s) is not an integer or out of range", string(args[3])))
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
	var nval float64
	var value float64
	if has {
		value, err = strconv.ParseFloat(string(kv.Value), 64)
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR old value(%s) is not an integer or out of range", string(kv.Value)))
			return
		}
		nval = value + valinc
		kv.Value = []byte(fmt.Sprintf("%f", nval))
		_, err = ms.eg.ID(kv.Id).Update(&kv)
	} else {
		nval = valinc
		kv.Value = []byte(fmt.Sprintf("%f", nval))
		_, err = ms.eg.Insert(&kv)
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteBulk([]byte(fmt.Sprintf("%f", nval)))
}
func (ms dbStorage) cmd_mget(conn redcon.Conn, args [][]byte) {
	if len(args) < 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkeys := []string{}
	for _, f := range args[1:] {
		rkeys = append(rkeys, string(f))
	}
	var kvs = []storage.Set{}
	err := ms.eg.In("rkey", rkeys).Cols("id", "rkey", "value").Find(&kvs)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteArray(len(args[1:]))
	for _, f := range args[1:] {
		p := false
		for _, kv := range kvs {
			if strings.EqualFold(string(f), kv.RKey) {
				p = true
				conn.WriteBulk(kv.Value)
			}
		}
		if !p {
			conn.WriteNull()
		}
	}
}
func (ms dbStorage) cmd_mset(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn < 3 || argn%2 == 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	kvs := []storage.Set{}
	for i := 1; i < argn; i = i + 2 {
		kvs = append(kvs, storage.Set{
			RKey:  string(args[i]),
			Value: args[i+1],
		})
	}
	var total_affected int64
	for _, kv := range kvs {
		val := kv.Value
		ok, err := ms.eg.Cols("id", "rkey", "value").Get(&kv)
		if err != nil {
			continue
		}
		var affected int64
		kv.Value = val
		if ok {
			affected, err = ms.eg.ID(kv.Id).Update(&kv)
		} else {
			affected, err = ms.eg.Insert(&kv)
		}
		if err != nil {
			continue
		}
		total_affected = total_affected + affected
	}
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_msetnx(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn < 3 || argn%2 == 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	kvs := []storage.Set{}
	for i := 1; i < argn; i = i + 2 {
		kvs = append(kvs, storage.Set{
			RKey:  string(args[i]),
			Value: args[i+1],
		})
	}
	var total_affected int64
	for _, kv := range kvs {
		val := kv.Value
		ok, err := ms.eg.Cols("id", "rkey", "value").Get(&kv)
		if err != nil {
			continue
		}
		var affected int64
		kv.Value = val
		if !ok {
			affected, err = ms.eg.Insert(&kv)
		}
		if err != nil {
			continue
		}
		total_affected = total_affected + affected
	}
	conn.WriteInt64(total_affected)
}

func (ms dbStorage) cmd_psetex(conn redcon.Conn, args [][]byte) {
	// TODO:
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
	has, err := ms.eg.Cols("id", "rkey", "value").Get(&kv)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	var affected int64
	if !has {
		kv.Value = args[2]
		affected, err = ms.eg.InsertOne(kv)
	} else {
		kv.Value = args[2]
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
