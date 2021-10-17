// Code generated by "stringer -type=EventType"; DO NOT EDIT.

package raft

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EventTypeHeartbeatReq-0]
	_ = x[EventTypeHeartbeatResp-1]
	_ = x[EventTypeClientReq-2]
	_ = x[EventTypeClientResp-3]
	_ = x[EventTypeVoteReq-4]
	_ = x[EventTypeVoteResp-5]
	_ = x[EventTypeAppendEntriesReq-6]
	_ = x[EventTypeAcceptEntriesResp-7]
	_ = x[EventTypeRefuseEntriesResp-8]
	_ = x[EventTypeInstallSnapReq-9]
}

const _EventType_name = "EventTypeHeartbeatReqEventTypeHeartbeatRespEventTypeClientReqMsgTypeClientRespEventTypeVoteReqEventTypeVoteRespEventTypeAppendEntriesReqEventTypeAcceptEntriesRespEventTypeRefuseEntriesRespEventTypeInstallSnapReq"

var _EventType_index = [...]uint8{0, 21, 43, 61, 78, 94, 111, 136, 162, 188, 211}

func (i EventType) String() string {
	if i < 0 || i >= EventType(len(_EventType_index)-1) {
		return "EventType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EventType_name[_EventType_index[i]:_EventType_index[i+1]]
}
