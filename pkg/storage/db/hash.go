package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dolfly/mydis/pkg/storage"
	"github.com/tidwall/redcon"
)

func (ms dbStorage) cmd_hdel(conn redcon.Conn, args [][]byte) {
	if len(args) < 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[1])
	fields := []string{}
	for _, f := range args[2:] {
		fields = append(fields, string(f))
	}
	var hashs = []storage.Hash{}
	err := ms.eg.Where("rkey = ?", rkey).In("field", fields).Unscoped().Find(&hashs)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if len(hashs) > 0 {
		var ints = []int64{}
		for _, h := range hashs {
			ints = append(ints, h.Id)
		}
		affected, err := ms.eg.In("id", ints).Unscoped().Delete(&storage.Hash{})
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
			return
		}
		conn.WriteInt64(affected)
	} else {
		conn.WriteInt64(0)
	}
}
func (ms dbStorage) cmd_hexists(conn redcon.Conn, args [][]byte) {
	if len(args) != 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	hash := storage.Hash{
		RKey:  string(args[1]),
		Field: string(args[2]),
	}
	ok, err := ms.eg.Where("rkey = ? and field = ?", hash.RKey, hash.Field).Exist(hash)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if ok {
		conn.WriteInt(1)
	} else {
		conn.WriteInt(0)
	}
}
func (ms dbStorage) cmd_hget(conn redcon.Conn, args [][]byte) {
	if len(args) != 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	hash := storage.Hash{
		RKey:  string(args[1]),
		Field: string(args[2]),
	}
	ok, err := ms.eg.Cols("id", "rkey", "field", "value").Get(&hash)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if ok {
		conn.WriteBulk(hash.Value)
	} else {
		conn.WriteNull()
	}
}
func (ms dbStorage) cmd_hgetall(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	hashs := []storage.Hash{}
	err := ms.eg.Where("rkey = ?", string(args[1])).Cols("id", "rkey", "field", "value").Find(&hashs)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteArray(len(hashs) * 2)
	for _, h := range hashs {
		conn.WriteString(h.Field)
		conn.WriteBulk(h.Value)
	}
}
func (ms dbStorage) cmd_hincrby(conn redcon.Conn, args [][]byte) {
	if len(args) != 4 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	valinc, err := strconv.ParseInt(string(args[3]), 10, 64)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR incr value(%s) is not an integer or out of range", string(args[3])))
		return
	}
	hash := storage.Hash{
		RKey:  string(args[1]),
		Field: string(args[2]),
	}
	ok, err := ms.eg.Cols("id", "rkey", "field", "value").Get(&hash)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	var nval int64
	var value int64
	if ok {
		value, err = strconv.ParseInt(string(hash.Value), 10, 64)
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR old value(%s) is not an integer or out of range", string(hash.Value)))
			return
		}
		nval = value + valinc
		hash.Value = []byte(fmt.Sprintf("%d", nval))
		_, err = ms.eg.ID(hash.Id).Update(&hash)
	} else {
		nval = valinc
		hash.Value = []byte(fmt.Sprintf("%d", nval))
		_, err = ms.eg.Insert(&hash)
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteInt64(nval)
}
func (ms dbStorage) cmd_hincrbyfloat(conn redcon.Conn, args [][]byte) {
	if len(args) != 4 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	valinc, err := strconv.ParseFloat(string(args[3]), 64)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR incr value(%s) is not an integer or out of range", string(args[3])))
		return
	}
	hash := storage.Hash{
		RKey:  string(args[1]),
		Field: string(args[2]),
	}
	ok, err := ms.eg.Cols("id", "rkey", "field", "value").Get(&hash)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	var nval float64
	var value float64
	if ok {
		value, err = strconv.ParseFloat(string(hash.Value), 64)
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR old value(%s) is not an integer or out of range", string(hash.Value)))
			return
		}
		nval = valinc + value
		hash.Value = []byte(fmt.Sprintf("%f", nval))
		_, err = ms.eg.ID(hash.Id).Update(&hash)
	} else {
		nval = valinc
		hash.Value = []byte(fmt.Sprintf("%f", nval))
		_, err = ms.eg.Insert(&hash)
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteBulk([]byte(fmt.Sprintf("%f", nval)))
}
func (ms dbStorage) cmd_hkeys(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[1])
	hashs := []storage.Hash{}
	err := ms.eg.Where("rkey = ?", rkey).Cols("field").Find(&hashs)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteArray(len(hashs))
	for _, h := range hashs {
		conn.WriteBulkString(h.Field)
	}
}
func (ms dbStorage) cmd_hlen(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[1])
	hashs := []storage.Hash{}
	err := ms.eg.Where("rkey = ?", rkey).Cols("field").Find(&hashs)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteInt(len(hashs))
}
func (ms dbStorage) cmd_hmget(conn redcon.Conn, args [][]byte) {
	if len(args) < 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[1])
	fields := []string{}
	for _, f := range args[2:] {
		fields = append(fields, string(f))
	}
	var hashs = []storage.Hash{}
	err := ms.eg.Where("rkey = ?", rkey).In("field", fields).Cols("id", "field", "value").Find(&hashs)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteArray(len(args[2:]))
	for _, f := range args[2:] {
		p := false
		for _, h := range hashs {
			if strings.EqualFold(string(f), h.Field) {
				p = true
				conn.WriteBulk(h.Value)
			}
		}
		if !p {
			conn.WriteNull()
		}
	}
}
func (ms dbStorage) cmd_hmset(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn < 4 || argn%2 != 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[1])
	hashs := []storage.Hash{}
	for i := 2; i < argn; i = i + 2 {
		hashs = append(hashs, storage.Hash{
			RKey:  rkey,
			Field: string(args[i]),
			Value: args[i+1],
		})
	}
	var total_affected int64
	for _, h := range hashs {
		val := h.Value
		ok, err := ms.eg.Where("rkey = ? and field = ?", h.RKey, h.Field).Get(&h)
		if err != nil {
			continue
		}
		var affected int64
		h.Value = val
		if ok {
			affected, err = ms.eg.ID(h.Id).Update(&h)
		} else {
			affected, err = ms.eg.Insert(&h)
		}
		if err != nil {
			continue
		}
		total_affected = total_affected + affected
	}
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_hscan(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_hset(conn redcon.Conn, args [][]byte) {
	if len(args) != 4 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	hash := storage.Hash{
		RKey:  string(args[1]),
		Field: string(args[2]),
	}
	ok, err := ms.eg.Cols("id", "rkey", "field").Get(&hash)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	hash.Value = args[3]
	var affected int64
	if ok {
		affected, err = ms.eg.ID(hash.Id).Cols("value").Update(&hash)
	} else {
		affected, err = ms.eg.Insert(&hash)
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteInt64(affected)
}
func (ms dbStorage) cmd_hsetnx(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
}
func (ms dbStorage) cmd_hstrlen(conn redcon.Conn, args [][]byte) {
	if len(args) != 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	hash := storage.Hash{
		RKey:  string(args[1]),
		Field: string(args[2]),
	}
	ok, err := ms.eg.Cols("id", "rkey", "field", "value").Get(&hash)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if ok {
		conn.WriteInt(len(hash.Value))
	} else {
		conn.WriteInt(0)
	}
}
func (ms dbStorage) cmd_hvals(conn redcon.Conn, args [][]byte) {
	if len(args) != 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	hashs := []storage.Hash{}
	err := ms.eg.Where("rkey = ?", string(args[1])).Cols("id", "rkey", "field", "value").Find(&hashs)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteArray(len(hashs))
	for _, h := range hashs {
		conn.WriteBulk(h.Value)
	}
}
