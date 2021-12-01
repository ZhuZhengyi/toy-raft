// message.go
package raft

import (
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

//MsgType msg event type
func (m *Message) MsgType() MsgType {
	if m.event != nil {
		return m.event.Type()
	}
	return MsgTypeUnkown
}

func (m *Message) From() uint64 {
	if addrPeer, ok := m.from.(*AddrPeer); ok {
		return addrPeer.peer
	}
	return 0
}

func (m *Message) To() uint64 {
	if addrPeer, ok := m.to.(*AddrPeer); ok {
		return addrPeer.peer
	}
	return 0
}

//String
func (m *Message) String() string {
	return fmt.Sprintf("{term:%v,%v,%v->%v}",
		m.term, m.event, m.from, m.to)
}

//Size message byte size with marshal
func (m *Message) Size() uint64 {
	return uint64(unsafe.Sizeof(m.From())) +
		uint64(unsafe.Sizeof(m.To())) +
		uint64(unsafe.Sizeof(m.term)) +
		uint64(unsafe.Sizeof(m.MsgType())) +
		m.event.Size()
}

//Marshal message to []byte
func (m *Message) Marshal(data []byte) {
	if len(data) < 28 {
		return
	}
	binary.BigEndian.PutUint64(data[0:], m.From())
	binary.BigEndian.PutUint64(data[8:], m.To())
	binary.BigEndian.PutUint64(data[16:], m.term)
	binary.BigEndian.PutUint32(data[24:], uint32(m.MsgType()))
	if len(data) > 28 {
		m.event.Marshal(data[28:])
	}
}

//Unmarshal message from []byte to message
func (m *Message) Unmarshal(data []byte) error {
	m.from = &AddrPeer{peer: binary.BigEndian.Uint64(data[0:8])}
	m.to = &AddrPeer{peer: binary.BigEndian.Uint64(data[8:16])}
	m.term = binary.BigEndian.Uint64(data[16:24])
	msgType := MsgType(binary.BigEndian.Uint32(data[24:28]))
	m.event = NewMsgEvent(msgType)

	return m.event.Unmarshal(data[28:])
}

//RecvFrom recv message from remote peer
func (m *Message) RecvFrom(r io.Reader) error {
	header := make([]byte, 8)
	_, err := r.Read(header[:8])
	if err != nil {
		logger.Error("read msg size error:%v", err)
		return err
	}
	size := binary.BigEndian.Uint64(header[:8])
	if size < 8 || size > 256*1024*1024 {
		logger.Debug("recv error msg size: %v \n", size)
		return ErrRecvMessage
	}

	data := make([]byte, size)
	_, err = r.Read(data[:size])
	if err != nil {
		logger.Error("read msg size error:%v", err)
		return err
	}
	m.Unmarshal(data[:size])

	return nil
}

//SendTo send  message to remote peer
func (m *Message) SendTo(w io.Writer) error {
	size := m.Size()
	if size <= 0 {
		return nil
	}

	msgSize := make([]byte, 8)
	binary.BigEndian.PutUint64(msgSize, size)
	if _, err := w.Write(msgSize[:8]); err != nil {
		logger.Fatal("Sent message(%v) err: %v\n", m, err)
		return err
	}

	data := make([]byte, size)
	m.Marshal(data)
	if _, err := w.Write(data[:size]); err != nil {
		logger.Error("send msg: %v err: %v", m, err)
		return err
	}

	return nil
}
