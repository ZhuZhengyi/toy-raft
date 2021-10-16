// message.go
package raft

type Message struct {
	from  Address
	to    Address
	term  uint64
	event MsgEvent
}

func NewMessage(from, to Address, term uint64, event MsgEvent) *Message {
	return &Message{
		from, to, term, event,
	}
}

func (m *Message) EventType() EventType {
	return m.event.Type()
}
