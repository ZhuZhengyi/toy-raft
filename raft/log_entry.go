// raft_log.go
package raft

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Entry struct {
	index   uint64
	term    uint64
	command []byte
}

func (e *Entry) String() string {
	return fmt.Sprintf("%v %v %v", e.index, e.term, e.command)
}

func (e *Entry) Size() uint64 {
	return uint64(8 + 8 + len(e.command))
}

func (e *Entry) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.index)
	binary.Write(buffer, binary.BigEndian, e.term)
	binary.Write(buffer, binary.BigEndian, e.command)

	return buffer.Bytes()
}

func (e *Entry) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, &e.index); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &e.term); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &e.command); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

type Entries []Entry

func (entries *Entries) Size() uint64 {
	s := uint64(0)
	for _, e := range []Entry(*entries) {
		s += e.Size()
	}
	return s
}

func (entries *Entries) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	for _, e := range []Entry(*entries) {
		binary.Write(buffer, binary.BigEndian, e.Marshal())
	}
	return buffer.Bytes()

}

func (entries *Entries) Unmarshal(data []byte) {

}
