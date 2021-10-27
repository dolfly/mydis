package storage

import (
	"time"

	"github.com/tidwall/redcon"
)

type Storage interface {
	Handler(conn redcon.Conn, cmd redcon.Command)
	Accept(conn redcon.Conn) bool
	Closed(conn redcon.Conn, err error)
}

type Namespace struct {
	Id      int64
	Name    string    `json:"name" xorm:"varchar(25) not null unique(ukey) 'name'"`
	RKey    string    `json:"rkey" xorm:"varchar(25) not null unique(ukey) 'rkey'"`
	Type    int       `json:"type" xorm:"tinyint 'type'"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
	Deleted time.Time `json:"deleted" xorm:"deleted"`
}

type KV struct {
	Id      int64
	RKey    string    `json:"rkey" xorm:"varchar(255) not null unique 'rkey'"`
	Value   []byte    `json:"value" xorm:"text 'value'"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
	Deleted time.Time `json:"deleted" xorm:"deleted"`
}
type Hash struct {
	Id      int64
	RKey    string    `json:"rkey" xorm:"varchar(255) not null unique(hkey) 'rkey'"`
	Field   string    `json:"hkey" xorm:"varchar(255) not null unique(hkey) 'field'"`
	Value   []byte    `json:"value" xorm:"text 'value'"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
	Deleted time.Time `json:"deleted" xorm:"deleted"`
}
type ZSet struct {
	RKey    string    `json:"rkey" xorm:"varchar(255) not null unique(zkey) 'rkey'"`
	Member  string    `json:"member" xorm:"varchar(255) not null unique(zkey) 'member'"`
	Score   int       `json:"score" xorm:"int 'score'"`
	Value   []byte    `json:"value" xorm:"text 'value'"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
	Deleted time.Time `json:"deleted" xorm:"deleted"`
}
