package raft

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
