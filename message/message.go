package message

import "io"

type Message interface {
	io.Closer
	ReadHeader(*Header) error
	ReadAuth(interface{}) error
	ReadBody(interface{}) error
	WriteMessage(header *Header, body interface{}) error
}

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

type NewMessageFunc func(closer io.ReadWriteCloser) Message

var MessageFuncMap map[Type]NewMessageFunc

func init() {
	MessageFuncMap = make(map[Type]NewMessageFunc)
	MessageFuncMap[GobType] = NewGobMessage
}
