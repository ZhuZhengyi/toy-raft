package raft

import (
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
