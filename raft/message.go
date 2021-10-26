// message.go
package raft

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

var (
	//ErrRecvMessage
	ErrRecvMessage = errors.New("recv message err")
)

//Message represents
type Message struct {
	from  Address
	to    Address
	term  uint64
	event MsgEvent
}

//NewMessage allocate a new message
func NewMessage(from, to Address, term uint64, event MsgEvent) *Message {
	return &Message{
		from, to, term, event,
	}
}

//EventType msg event type
func (m *Message) EventType() EventType {
	if m.event != nil {
		return m.event.Type()
	}
	return EventTypeUnkown
}

//String
func (m *Message) String() string {
	return fmt.Sprintf("Message(%v -> %v): term(%v) event(%v)",
		m.from, m.to, m.term, m.event)
}

//Size message byte size with marshal
func (m *Message) Size() uint64 {
	return uint64(unsafe.Sizeof(m.term)) +
		uint64(unsafe.Sizeof(m.EventType())) +
		m.event.Size()
}

//Marshal message to []byte
func (m *Message) Marshal() []byte {
	datas := takeBytes()
	defer putBytes(datas)

	buffer := bytes.NewBuffer(datas)
	binary.Write(buffer, binary.BigEndian, m.term)
	binary.Write(buffer, binary.BigEndian, m.EventType())
	binary.Write(buffer, binary.BigEndian, m.event.Marshal())
	return buffer.Bytes()
}

//Unmarshal message from []byte to message
func (m *Message) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, &m.term); err != nil {
		logger.Warn("unmarshal %v error: %v", m, err)
		return err
	}

	var msgEventType EventType
	binary.Read(buffer, binary.BigEndian, &msgEventType)
	m.event = NewMsgEvent(msgEventType)

	return m.event.Unmarshal(data[8:])
}

//RecvFrom recv message from remote peer
func (m *Message) RecvFrom(r io.Reader) error {
	buf := takeBytes()
	defer putBytes(buf)

	_, err := r.Read(buf[:8])
	if err != nil {
		logger.Error("read msg size error:%v", err)
		return err
	}
	size := binary.BigEndian.Uint64(buf[:8])
	if size < 8 || size > 256*1024*1024 {
		logger.Error("recv error msg size: %v \n", size)
		return ErrRecvMessage
	}

	_, err = r.Read(buf[:size])
	if err != nil {
		logger.Error("read msg size error:%v", err)
		return err
	}
	//msg := takeMessage()
	m.Unmarshal(buf[:size])

	return nil
}

//SendTo send  message to remote peer
func (m *Message) SendTo(w io.Writer) error {
	size := m.Size()
	bytes := takeBytes()
	binary.BigEndian.PutUint64(bytes, size)
	if _, err := w.Write(bytes[:8]); err != nil {
		logger.Fatal("Sent message(%v) err: %v\n", m, err)
		return err
	}

	data := m.Marshal()
	if _, err := w.Write(data); err != nil {
		logger.Error("send msg: %v err: %v", m, err)
		return err
	}

	return nil
}
