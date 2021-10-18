// message.go
package raft

import (
	"fmt"
	"io"
)

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

func (m *Message) Size() uint64 {
	return 0
}

func (m *Message) Marshal() []byte {
	return nil
}

func (m *Message) Unmarshal(msgData []byte) {
}

func (m *Message) Recv(r io.Reader) error {

	return nil
}

func (m *Message) Send(w io.Writer) error {

	data := m.Marshal()

	if _, err := w.Write(data); err != nil {
		logger.Error("send msg: %v err: %v", m, err)
		return err
	}

	return nil
}
