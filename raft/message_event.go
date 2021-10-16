package raft

import (
	"github.com/google/uuid"
)

type ReqId = uuid.UUID
type MsgType int32

const (
	MsgTypeUnkown       MsgType = -1
	MsgTypeHeartbeatReq MsgType = iota
	MsgTypeHeartbeatResp
	MsgTypeClientReq
	MsgTypeClientResp
	MsgTypeVoteReq
	MsgTypeVoteResp
	MsgTypeAppendEntriesReq
	MsgTypeAcceptEntriesResp
	MsgTypeRefuseEntriesResp
	MsgTypeInstallSnapReq
)

type MsgEvent interface {
	Type() MsgType
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

func (e *EventHeartbeatReq) Type() MsgType {
	return MsgTypeHeartbeatReq
}

func (e *EventHeartbeatResp) Type() MsgType {
	return MsgTypeHeartbeatResp
}

func (e *EventClientReq) Type() MsgType {
	return MsgTypeClientReq
}

func (e *EventClientResp) Type() MsgType {
	return MsgTypeClientResp
}

func (e *EventSolicitVoteReq) Type() MsgType {
	return MsgTypeVoteReq
}

func (e *EventGrantVoteResp) Type() MsgType {
	return MsgTypeVoteResp
}

func (e *EventAppendEntriesReq) Type() MsgType {
	return MsgTypeAppendEntriesReq
}

func (e *EventAcceptEntriesResp) Type() MsgType {
	return MsgTypeAcceptEntriesResp
}

func (e *EventRefuseEntriesResp) Type() MsgType {
	return MsgTypeRefuseEntriesResp
}
