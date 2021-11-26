package raft

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/google/uuid"
)

//go:generate stringer -type=MsgEventType  -linecomment

type ReqId = uuid.UUID
type MsgEventType int32

const (
	EventTypeUnkown            MsgEventType = -1   // EventUknown
	EventTypeHeartbeatReq      MsgEventType = iota // EventHeartbeatReq
	EventTypeHeartbeatResp                         // EventHeartbeatResponse
	EventTypeClientReq                             // EventClientReq
	EventTypeClientResp                            // EventClientResponse
	EventTypeVoteReq                               // EventVoteReq
	EventTypeVoteResp                              // EventVoteResponse
	EventTypeAppendEntriesReq                      // EventAppendEntriesReq
	EventTypeAcceptEntriesResp                     // EventAcceptEntriesResponse
	EventTypeRefuseEntriesResp                     // EventRefuseEntriesResponse
	EventTypeInstallSnapReq                        // EventInstallSnapshot
)

//MsgEvent message event
type MsgEvent interface {
	Type() MsgEventType
	Size() uint64
	Marshal() []byte
	Unmarshal(data []byte) error
}

type queuedEvent struct {
	from  Address
	event MsgEvent
}

func (e *queuedEvent) EventType() MsgEventType {
	if e.event != nil {
		return e.event.Type()

	}
	return EventTypeUnkown

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
func NewMsgEvent(eventType MsgEventType) MsgEvent {
	switch eventType {
	case EventTypeHeartbeatReq:
		return new(EventHeartbeatReq)
	case EventTypeHeartbeatResp:
		return new(EventHeartbeatResp)
	case EventTypeVoteReq:
		return new(EventSolicitVoteReq)
	case EventTypeVoteResp:
		return new(EventGrantVoteResp)
	case EventTypeClientReq:
		return new(EventClientResp)
	case EventTypeAppendEntriesReq:
		return new(EventAppendEntriesReq)
	case EventTypeAcceptEntriesResp:
		return new(EventAcceptEntriesResp)
	case EventTypeRefuseEntriesResp:
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