package raft

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/google/uuid"
)

//go:generate stringer -type=EventType  -linecomment

type ReqId = uuid.UUID
type EventType int32

const (
	EventTypeHeartbeatReq      EventType = iota // EventHeartbeatReq
	EventTypeHeartbeatResp                      // EventHeartbeatResponse
	EventTypeClientReq                          // EventClientReq
	EventTypeClientResp                         // EventClientResponse
	EventTypeVoteReq                            // EventVoteReq
	EventTypeVoteResp                           // EventVoteResponse
	EventTypeAppendEntriesReq                   // EventAppendEntriesReq
	EventTypeAcceptEntriesResp                  // EventAcceptEntriesResponse
	EventTypeRefuseEntriesResp                  // EventRefuseEntriesResponse
	EventTypeInstallSnapReq                     // EventInstallSnapshot
)

type MsgEvent interface {
	Type() EventType
	Size() uint64
	Marshal() []byte
	Unmarshal(data []byte)
}

type queuedEvent struct {
	from  Address
	event MsgEvent
}

var (
	NOOPCommand = []byte{}
)

var (
	_ (MsgEvent) = (*EventHeartbeatReq)(nil)
	_ (MsgEvent) = (*EventHeartbeatResp)(nil)
	_ (MsgEvent) = (*EventSolicitVoteReq)(nil)
	_ (MsgEvent) = (*EventGrantVoteResp)(nil)
	_ (MsgEvent) = (*EventClientReq)(nil)
	_ (MsgEvent) = (*EventClientResp)(nil)
	_ (MsgEvent) = (*EventAppendEntriesReq)(nil)
)

type EventHeartbeatReq struct {
	commitIndex uint64
	commitTerm  uint64
}

type EventHeartbeatResp struct {
	commitIndex  uint64
	hasCommitted bool
}

type EventSolicitVoteReq struct {
	lastIndex uint64
	lastTerm  uint64
}

type EventGrantVoteResp struct {
}

type EventClientReq struct {
	id      ReqId
	request Request
}

func NewEventClientReq(req Request) *EventClientReq {
	uuid, _ := uuid.NewRandom()
	return &EventClientReq{
		id:      uuid,
		request: req,
	}
}

type EventClientResp struct {
	id       ReqId
	response Response
}

type EventAppendEntriesReq struct {
	baseIndex uint64
	baseTerm  uint64
	entries   []Entry
}

type EventAcceptEntriesResp struct {
}

type EventRefuseEntriesResp struct {
}

func (e *EventHeartbeatReq) Size() uint64 {
	return uint64(unsafe.Sizeof(e))
}

func (a *EventHeartbeatResp) Size() uint64 {
	return uint64(unsafe.Sizeof(a))
}

func (e *EventClientReq) Size() uint64 {
	return uint64(unsafe.Sizeof(e))
}

func (a *EventClientResp) Size() uint64 {
	return uint64(unsafe.Sizeof(a))
}

func (a *EventSolicitVoteReq) Size() uint64 {
	return uint64(unsafe.Sizeof(a))
}

func (a *EventGrantVoteResp) Size() uint64 {
	return uint64(unsafe.Sizeof(a))
}

func (a *EventAppendEntriesReq) Size() uint64 {
	return uint64(unsafe.Sizeof(a))
}

func (a *EventAcceptEntriesResp) Size() uint64 {
	return uint64(unsafe.Sizeof(a))
}

func (a *EventRefuseEntriesResp) Size() uint64 {
	return uint64(unsafe.Sizeof(a))
}

func (e *EventHeartbeatReq) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.commitIndex)
	binary.Write(buffer, binary.BigEndian, e.commitTerm)

	return buffer.Bytes()
}

func (e *EventHeartbeatReq) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventHeartbeatResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.commitIndex)
	binary.Write(buffer, binary.BigEndian, e.hasCommitted)

	return buffer.Bytes()
}

func (e *EventHeartbeatResp) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventClientReq) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})

	return buffer.Bytes()
}

func (e *EventClientReq) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventClientResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})

	return buffer.Bytes()
}

func (e *EventClientResp) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventSolicitVoteReq) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.lastIndex)
	binary.Write(buffer, binary.BigEndian, e.lastTerm)

	return buffer.Bytes()
}

func (e *EventSolicitVoteReq) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventGrantVoteResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})

	return buffer.Bytes()
}

func (e *EventGrantVoteResp) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventAppendEntriesReq) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.baseIndex)
	binary.Write(buffer, binary.BigEndian, e.baseTerm)
	for _, entry := range e.entries {
		binary.Write(buffer, binary.BigEndian, entry)
	}

	return buffer.Bytes()
}

func (e *EventAppendEntriesReq) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventAcceptEntriesResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.commitIndex)
	binary.Write(buffer, binary.BigEndian, e.hasCommitted)

	return buffer.Bytes()
}

func (e *EventAcceptEntriesResp) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}

func (e *EventRefuseEntriesResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.commitIndex)
	binary.Write(buffer, binary.BigEndian, e.hasCommitted)

	return buffer.Bytes()
}

func (e *EventRefuseEntriesResp) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
	}
}
