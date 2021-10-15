package raft

import (
	"github.com/google/uuid"
)

type ReqId = uuid.UUID
type MsgType uint32

const (
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
)

type MsgEvent interface {
	Type() MsgType
}

type EventHeartbeat struct {
	commitIndex uint64
	commitTerm  uint64
}

func (e *EventHeartbeat) Type() MsgType {
	return MsgHeartbeat
}

type EventHeartbeatResp struct {
	commitIndex  uint64
	hasCommitted bool
}

func (e *EventHeartbeatResp) Type() MsgType {
	return MsgHeartbeatResp
}

type EventSolicitVoteReq struct {
	lastIndex uint64
	lastTerm  uint64
}

func (e *EventSolicitVoteReq) Type() MsgType {
	return MsgSolictVoteReq
}

type EventGrantVoteResp struct {
}

func (e *EventGrantVoteResp) Type() MsgType {
	return MsgGrantVoteResp
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

func (e *EventClientReq) Type() MsgType {
	return MsgClientReq
}

type EventClientResp struct {
	id       ReqId
	response Response
}

func (e *EventClientResp) Type() MsgType {
	return MsgClientResp
}
