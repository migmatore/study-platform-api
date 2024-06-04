package websocket

import "github.com/migmatore/study-platform-api/internal/core"

type MessageType int

const (
	AuthRequest MessageType = iota + 1
	VirtualPointer
	Call
	NewRoom
	ErrorResp
)

type ErrorType int

const (
	ExpiredTokenError ErrorType = iota
)

type Receiver struct {
	Id   int
	role core.RoleType
}

type Message struct {
	MsgType MessageType `json:"msg_type"`
	Data    []byte      `json:"data"`
	To      []Receiver
}

func NewMessage(data []byte, to []Receiver) *Message {
	return &Message{Data: data, To: to}
}
