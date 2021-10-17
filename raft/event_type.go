package raft

func (e *EventHeartbeatReq) Type() EventType {
	return EventTypeHeartbeatReq
}

func (e *EventHeartbeatResp) Type() EventType {
	return EventTypeHeartbeatResp
}

func (e *EventClientReq) Type() EventType {
	return EventTypeClientReq
}

func (e *EventClientResp) Type() EventType {
	return EventTypeClientResp
}

func (e *EventSolicitVoteReq) Type() EventType {
	return EventTypeVoteReq
}

func (e *EventGrantVoteResp) Type() EventType {
	return EventTypeVoteResp
}

func (e *EventAppendEntriesReq) Type() EventType {
	return EventTypeAppendEntriesReq
}

func (e *EventAcceptEntriesResp) Type() EventType {
	return EventTypeAcceptEntriesResp
}

func (e *EventRefuseEntriesResp) Type() EventType {
	return EventTypeRefuseEntriesResp
}
