package raft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/google/uuid"
)

//go:generate stringer -type=MsgType  -linecomment

type ReqId = uuid.UUID
type MsgType int32

const (
	MsgTypeUnkown            MsgType = -1   // EventUknown
	MsgTypeHeartbeatReq      MsgType = iota // EventHeartbeatReq
	MsgTypeHeartbeatResp                    // EventHeartbeatResponse
	MsgTypeClientReq                        // EventClientReq
	MsgTypeClientResp                       // EventClientResponse
	MsgTypeVoteReq                          // EventVoteReq
	MsgTypeVoteResp                         // EventVoteResponse
	MsgTypeAppendEntriesReq                 // EventAppendEntriesReq
	MsgTypeAcceptEntriesResp                // EventAcceptEntriesResponse
	MsgTypeRefuseEntriesResp                // EventRefuseEntriesResponse
	MsgTypeInstallSnapReq                   // EventInstallSnapshot
)

//MsgEvent message event
type MsgEvent interface {
	Type() MsgType
	Size() uint64
	String() string
	Marshal() []byte
	Unmarshal(data []byte) error
}

type queuedEvent struct {
	from  Address
	event MsgEvent
}

func (e *queuedEvent) MsgType() MsgType {
	if e.event != nil {
		return e.event.Type()

	}
	return MsgTypeUnkown

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

//NewMsgEvent allocate a new MsgEvent by eventType
func NewMsgEvent(eventType MsgType) MsgEvent {
	switch eventType {
	case MsgTypeHeartbeatReq:
		return new(EventHeartbeatReq)
	case MsgTypeHeartbeatResp:
		return new(EventHeartbeatResp)
	case MsgTypeVoteReq:
		return new(EventSolicitVoteReq)
	case MsgTypeVoteResp:
		return new(EventGrantVoteResp)
	case MsgTypeClientReq:
		return new(EventClientResp)
	case MsgTypeAppendEntriesReq:
		return new(EventAppendEntriesReq)
	case MsgTypeAcceptEntriesResp:
		return new(EventAcceptEntriesResp)
	case MsgTypeRefuseEntriesResp:
		return new(EventRefuseEntriesResp)
	default:
		return new(EventHeartbeatReq)
	}
}

//EventHeartbeatReq event heartbeat req from leader -> follower
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

//
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

func (e *EventHeartbeatReq) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, &e.commitIndex); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	if err := binary.Read(buffer, binary.BigEndian, &e.commitTerm); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventHeartbeatResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.commitIndex)
	binary.Write(buffer, binary.BigEndian, e.hasCommitted)

	return buffer.Bytes()
}

func (e *EventHeartbeatResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, &e.commitIndex); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	if err := binary.Read(buffer, binary.BigEndian, &e.hasCommitted); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventClientReq) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})

	return buffer.Bytes()
}

func (e *EventClientReq) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventClientResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})

	return buffer.Bytes()
}

func (e *EventClientResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventSolicitVoteReq) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, e.lastIndex)
	binary.Write(buffer, binary.BigEndian, e.lastTerm)

	return buffer.Bytes()
}

func (e *EventSolicitVoteReq) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventGrantVoteResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})

	return buffer.Bytes()
}

func (e *EventGrantVoteResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}

	return nil
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

func (e *EventAppendEntriesReq) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventAcceptEntriesResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	return buffer.Bytes()
}

func (e *EventAcceptEntriesResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventRefuseEntriesResp) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})

	return buffer.Bytes()
}

func (e *EventRefuseEntriesResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventHeartbeatReq) String() string {
	return fmt.Sprintf("{commitIndex: %v, commitTerm: %v}", e.commitIndex, e.commitTerm)
}

func (e *EventHeartbeatResp) String() string {
	return fmt.Sprintf("{commitIndex: %v, hasCommitted: %v}", e.commitIndex, e.hasCommitted)
}

func (e *EventSolicitVoteReq) String() string {
	return fmt.Sprintf("{lastIndex: %v, lastTerm: %v}", e.lastIndex, e.lastTerm)
}

func (e *EventGrantVoteResp) String() string {
	return "{}"
}

func (e *EventClientReq) String() string {
	return fmt.Sprintf("{id: %v, request: %v}", e.id, e.request)
}

func (e *EventClientResp) String() string {
	return fmt.Sprintf("{id: %v, response: %v}", e.id, e.response)
}

func (e *EventAppendEntriesReq) String() string {
	return fmt.Sprintf("{baseIndex: %v, baseTerm: %v}", e.baseIndex, e.baseTerm)
}
func (e *EventAcceptEntriesResp) String() string {
	return "{}"
}
func (e *EventRefuseEntriesResp) String() string {
	return "{}"
}
