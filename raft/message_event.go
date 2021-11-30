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
	Marshal(data []byte)
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

func (e *EventHeartbeatReq) Marshal(data []byte) {
	binary.BigEndian.PutUint64(data[0:], e.commitIndex)
	if len(data) > 8 {
		binary.BigEndian.PutUint64(data[8:], e.commitTerm)
	}
}

func (e *EventHeartbeatReq) Unmarshal(data []byte) error {
	e.commitIndex = binary.BigEndian.Uint64(data[:8])
	if len(data) >= 16 {
		e.commitTerm = binary.BigEndian.Uint64(data[8:16])
	}
	return nil
}

func (e *EventHeartbeatResp) Marshal(data []byte) {
	binary.BigEndian.PutUint64(data[0:], e.commitIndex)
	//binary.BigEndian.PutUint16(data[8:], e.hasCommitted)
}

func (e *EventHeartbeatResp) Unmarshal(data []byte) error {
	e.commitIndex = binary.BigEndian.Uint64(data[:8])
	//e.hasCommitted = binary.BigEndian.Uint64(data[8:16])
	return nil
}

func (e *EventClientReq) Marshal(data []byte) {
	//TODO:

}

func (e *EventClientReq) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventClientResp) Marshal(data []byte) {
	//TODO:
}

func (e *EventClientResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventSolicitVoteReq) Marshal(data []byte) {
	binary.BigEndian.PutUint64(data[0:], e.lastIndex)
	if len(data) > 8 {
		binary.BigEndian.PutUint64(data[8:], e.lastTerm)
	}
}

func (e *EventSolicitVoteReq) Unmarshal(data []byte) error {
	e.lastIndex = binary.BigEndian.Uint64(data[:8])
	if len(data) >= 16 {
		e.lastTerm = binary.BigEndian.Uint64(data[8:16])
	}
	return nil
}

func (e *EventGrantVoteResp) Marshal(data []byte) {
}

func (e *EventGrantVoteResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}

	return nil
}

func (e *EventAppendEntriesReq) Marshal(data []byte) {
	binary.BigEndian.PutUint64(data[0:], e.baseIndex)
	binary.BigEndian.PutUint64(data[8:], e.baseTerm)
	//TODO:
	//for _, entry := range e.entries {
	//    binary.Write(buffer, binary.BigEndian, entry)
	//}

	//return buffer.Bytes()
}

func (e *EventAppendEntriesReq) Unmarshal(data []byte) error {
	//TODO:
	//buffer := bytes.NewBuffer(data)

	//if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
	//    logger.Warn("unmarshal %v error: %v", e, err)
	//    return err
	//}
	return nil
}

func (e *EventAcceptEntriesResp) Marshal(data []byte) {
	//TODO:
}

func (e *EventAcceptEntriesResp) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, e); err != nil {
		logger.Warn("unmarshal %v error: %v", e, err)
		return err
	}
	return nil
}

func (e *EventRefuseEntriesResp) Marshal(data []byte) {

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
	return fmt.Sprintf("%v{commitIndex: %v, commitTerm: %v}", e.Type(), e.commitIndex, e.commitTerm)
}

func (e *EventHeartbeatResp) String() string {
	return fmt.Sprintf("%v{commitIndex: %v, hasCommitted: %v}", e.Type(), e.commitIndex, e.hasCommitted)
}

func (e *EventSolicitVoteReq) String() string {
	return fmt.Sprintf("%v{lastIndex: %v, lastTerm: %v}", e.Type(), e.lastIndex, e.lastTerm)
}

func (e *EventGrantVoteResp) String() string {
	return fmt.Sprintf("%v{}", e.Type())
}

func (e *EventClientReq) String() string {
	return fmt.Sprintf("%v{id: %v, request: %v}", e.Type(), e.id, e.request)
}

func (e *EventClientResp) String() string {
	return fmt.Sprintf("%v{id: %v, response: %v}", e.Type(), e.id, e.response)
}

func (e *EventAppendEntriesReq) String() string {
	return fmt.Sprintf("%v{baseIndex: %v, baseTerm: %v}", e.Type(), e.baseIndex, e.baseTerm)
}
func (e *EventAcceptEntriesResp) String() string {
	return fmt.Sprintf("%v{}", e.Type())
}
func (e *EventRefuseEntriesResp) String() string {
	return fmt.Sprintf("%v{}", e.Type())
}
