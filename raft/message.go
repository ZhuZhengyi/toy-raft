// message.go
package raft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
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
	return m.from.Size() +
		m.to.Size() +
		uint64(unsafe.Sizeof(m.term)) +
		m.event.Size()
}

func (m *Message) Marshal() []byte {
	datas := takeBytes()
	defer putBytes(datas)

	buffer := bytes.NewBuffer(datas)
	binary.Write(buffer, binary.BigEndian, m.term)
	binary.Write(buffer, binary.BigEndian, m.event.Marshal())
	return buffer.Bytes()
}

func (m *Message) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, m.term); err != nil {
		logger.Warn("unmarshal %v error: %v", m, err)
	}

	m.event.Unmarshal(data[8:])
}

func (m *Message) RecvFrom(r io.Reader) error {
	buf := takeBytes()
	defer putBytes(buf)

	r.Read(buf[:8])

	//r.Read(p []byte)

	return nil
}

func (m *Message) SendTo(w io.Writer) error {
	size := m.Size()
	bytes := takeBytes()
	binary.BigEndian.PutUint64(bytes, size)
	if _, err := w.Write(bytes[:8]); err != nil {
	}

	data := m.Marshal()
	if _, err := w.Write(data); err != nil {
		logger.Error("send msg: %v err: %v", m, err)
		return err
	}

	return nil
}
