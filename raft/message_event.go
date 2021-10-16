package raft

import (
	"github.com/google/uuid"
)

type ReqId = uuid.UUID
type MsgType int32

const (
	MsgUnkown MsgType = -1

	MsgHeartbeat MsgType = iota
	MsgHeartbeatResp
	MsgClientReq
	MsgClientResp
	MsgSolictVoteReq
	MsgGrantVoteResp
	MsgAppendEntriesReq
	MsgAcceptEntriesResp
	MsgRefuseEntriesResp
	MsgInstallSnapReq

	MsgNameUnkown        = "MsgUnkown"
	MsgNameHeartbeat     = "MsgHeartbeatReq"
	MsgNameHeartbeatResp = "MsgHeartbeatResp"
	MsgNameClientReq     = "MsgClientReq"
	MsgNameClientResp
	MsgNameSolictVoteReq
	MsgNameGrantVoteResp
	MsgNameAppendEntriesReq
	MsgNameAcceptEntriesResp
	MsgNameRefuseEntriesResp
	MsgNameInstallSnapReq
)

type MsgEvent interface {
	Type() MsgType
	String() string
}

type msgEvent struct {
}

func (e *msgEvent) Type() MsgType {
	switch MsgEvent(e).(type) {
	case *EventHeartbeatReq:
		return MsgHeartbeat
	case *EventHeartbeatResp:
		return MsgHeartbeatResp
	case *EventSolicitVoteReq:
		return MsgSolictVoteReq
	case *EventGrantVoteResp:
		return MsgGrantVoteResp
	case *EventAppendEntriesReq:
		return MsgAppendEntriesReq
	case *EventClientReq:
		return MsgClientReq
	case *EventClientResp:
		return MsgClientResp
	}

	return MsgUnkown
}

func (e *msgEvent) String() string {
	switch MsgEvent(e).(type) {
	case *EventHeartbeatReq:
		return MsgNameHeartbeat
	case *EventHeartbeatResp:
		return MsgNameHeartbeatResp
	case *EventSolicitVoteReq:
		return MsgNameSolictVoteReq
	case *EventGrantVoteResp:
		return MsgNameGrantVoteResp
	case *EventAppendEntriesReq:
		return MsgNameAppendEntriesReq
	case *EventClientReq:
		return MsgNameClientReq
	case *EventClientResp:
		return MsgNameClientResp
	}

	return MsgNameUnkown
}

type EventHeartbeatReq struct {
	msgEvent
	commitIndex uint64
	commitTerm  uint64
}

type EventHeartbeatResp struct {
	msgEvent
	commitIndex  uint64
	hasCommitted bool
}

type EventSolicitVoteReq struct {
	msgEvent
	lastIndex uint64
	lastTerm  uint64
}

type EventGrantVoteResp struct {
	msgEvent
}

type EventClientReq struct {
	msgEvent
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
	msgEvent
	id       ReqId
	response Response
}

type EventAppendEntriesReq struct {
	msgEvent
	baseIndex uint64
	baseTerm  uint64
	entries   []Entry
}
