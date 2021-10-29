package db

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func (ms dbStorage) cmd_ping(conn redcon.Conn, args [][]byte) {
	conn.WriteString("PONG")
}
func (ms dbStorage) cmd_quit(conn redcon.Conn, args [][]byte) {
	conn.WriteString("OK")
	conn.Close()
}
func (ms dbStorage) cmd_dblist(conn redcon.Conn, args [][]byte) {
	//conn.WriteString("OK")
}
func (ms dbStorage) cmd_dbtype(conn redcon.Conn, args [][]byte) {
	conn.WriteArray(len(dict) * 2)
	for k, v := range dict {
		conn.WriteBulkString(k)
		conn.WriteBulkString(string(v))
	}
}
func (ms dbStorage) cmd_dbdump(conn redcon.Conn, args [][]byte) {
	argn := len(args)
	if argn < 2 {
		conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", string(args[0])))
		return
	}
	var err error
	if argn >= 3 {
		if tp, ok := dict[string(args[2])]; ok {
			err = ms.eg.DumpAllToFile(string(args[1]), tp)
		} else {
			err = ms.eg.DumpAllToFile(string(args[1]))
		}
	} else {
		err = ms.eg.DumpAllToFile(string(args[1]))
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %s", err.Error()))
		return
	}
	conn.WriteString("OK")
}
