// message.go
package raft

import "fmt"

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

func (m *Message) String() string {
	return fmt.Sprintf("Message(%v -> %v): term(%v) event(%v)",
		m.from, m.to, m.term, m.event)
}
