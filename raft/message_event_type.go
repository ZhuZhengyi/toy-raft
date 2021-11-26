package raft

func (e *EventHeartbeatReq) Type() MsgEventType {
	return EventTypeHeartbeatReq
}

func (e *EventHeartbeatResp) Type() MsgEventType {
	return EventTypeHeartbeatResp
}

func (e *EventClientReq) Type() MsgEventType {
	return EventTypeClientReq
}

func (e *EventClientResp) Type() MsgEventType {
	return EventTypeClientResp
}

func (e *EventSolicitVoteReq) Type() MsgEventType {
	return EventTypeVoteReq
}

func (e *EventGrantVoteResp) Type() MsgEventType {
	return EventTypeVoteResp
}

func (e *EventAppendEntriesReq) Type() MsgEventType {
	return EventTypeAppendEntriesReq
}

func (e *EventAcceptEntriesResp) Type() MsgEventType {
	return EventTypeAcceptEntriesResp
}

func (e *EventRefuseEntriesResp) Type() MsgEventType {
	return EventTypeRefuseEntriesResp
}
