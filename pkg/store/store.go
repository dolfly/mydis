package store

import (
	"time"

	"github.com/tidwall/redcon"
)

type Store interface {
	Handler(conn redcon.Conn, cmd redcon.Command)
	Accept(conn redcon.Conn) bool
	Closed(conn redcon.Conn, err error)
}

type Namespace struct {
	Name string `json:"ns"`
	RKey string `json:"rkey"`
	Type string `json:"type"`
}

type KV struct {
	RKey    string    `json:"rkey"`
	Value   string    `json:"value"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Deleted time.Time `json:"deleted"`
}
type Hash struct {
	RKey    string    `json:"rkey"`
	HKey    string    `json:"hkey"`
	Value   string    `json:"value"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Deleted time.Time `json:"deleted"`
}
type ZSet struct {
	RKey   string `json:"rkey"`
	Member string `json:"member"`
	Score  int    `json:"score"`
	Value  string `json:"value"`
}
