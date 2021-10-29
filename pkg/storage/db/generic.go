package db

import (
	"github.com/tidwall/redcon"
)

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
