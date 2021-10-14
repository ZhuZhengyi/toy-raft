// message.go
package raft

type Address = uint32

const (
	AddrClient Address = iota
	AddrLocal
	AddrPeer
	AddrPeers
)

//type Message interface {
//    Type() MsgType
//    From() Address
//    To() Address
//    //Term() uint64
//}

type Message struct {
	from  Address
	to    Address
	term  uint64
	event MsgEvent
}

func NewMessage(from, to Address, term uint64, event MsgEvent) *Message {
	return &Message{}
}

func (m *Message) Type() MsgType {
	return m.event.Type()
}
