package db

import (
	"fmt"
	"strings"

	"github.com/dolfly/mydis/pkg/storage"
	"github.com/tidwall/redcon"
)

func (ms dbStorage) cmd_blpop(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_brpop(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_brpoplpush(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_lindex(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_linsert(conn redcon.Conn, args [][]byte) {
	if len(args) != 5 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[1])
	pos := string(args[2])
	pivot := args[3]
	value := args[4]

	switch strings.ToLower(pos) {
	case "BEFORE":
	case "AFTER":
	default:
		conn.WriteError("ERR syntax error")
		return
	}

	plst := storage.List{
		RKey:  rkey,
		Value: pivot,
	}
	has, err := ms.eg.Asc("index").Cols("id", "rkey", "index", "value").Get(&plst)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	sess := ms.eg.NewSession()
	defer sess.Close()

	if has {
		nlst := storage.List{
			RKey: rkey,
			//Index: plst.Index + 1,
			Value: value,
		}
		if err := sess.Begin(); err == nil {
			if _, err = sess.Table(new(storage.List)).Where("rkey = ? and index > ?", rkey, 0).Cols("index").Update(map[string]interface{}{
				"index": "index + 1",
			}); err != nil {
				conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
				return
			}
			if _, err = sess.Insert(&nlst); err != nil {
				conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
				return
			}
			err = sess.Commit()
			if err != nil {
				conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
				return
			}
			cnt, err := ms.eg.Where("rkey = ?", rkey).Count(new(storage.List))
			if err != nil {
				conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
				return
			}
			conn.WriteInt64(cnt)
		} else {
			conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
			return
		}
	} else {
		conn.WriteInt64(0)
	}
}
func (ms dbStorage) cmd_llen(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_lpop(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_lpush(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn < 3 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[1])

	plst := storage.List{
		RKey: rkey,
	}
	lists := []storage.List{}
	//index := 0
	for i := argn; i > 2; i-- {
		lists = append(lists, storage.List{
			RKey: rkey,
			//Index: int64(index + argn - i),
			Value: args[i-1],
		})
	}

	has, err := ms.eg.Asc("index").Cols("id", "rkey", "index", "value").Get(&plst)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	if has {
		sess := ms.eg.NewSession()
		if err := sess.Begin(); err != nil {
			conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
			return
		}
		affected, err := sess.Table(new(storage.List)).Where("rkey = ?", rkey).Update(map[string]interface{}{
			"index": fmt.Sprintf("index + %d", argn-2),
		})
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
			return
		}
		_, err = sess.Insert(&lists)
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
			return
		}
		err = sess.Commit()
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
			return
		}
		conn.WriteInt(int(affected) + argn - 2)
	} else {
		_, err = ms.eg.Insert(&lists)
		if err != nil {
			conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
			return
		}
		conn.WriteInt(argn - 2)
	}

}
func (ms dbStorage) cmd_lpushx(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_lrange(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_lrem(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_lset(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_ltrim(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_rpop(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_rpoplpush(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_rpush(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
func (ms dbStorage) cmd_rpushx(conn redcon.Conn, args [][]byte) {
	if len(args) < 0 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	rkey := string(args[0])
	_ = rkey

}
