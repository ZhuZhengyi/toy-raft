// Code generated by "stringer -type=MsgType -linecomment"; DO NOT EDIT.

package raft

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MsgTypeUnkown - -1]
	_ = x[MsgTypeHeartbeatReq-1]
	_ = x[MsgTypeHeartbeatResp-2]
	_ = x[MsgTypeClientReq-3]
	_ = x[MsgTypeClientResp-4]
	_ = x[MsgTypeVoteReq-5]
	_ = x[MsgTypeVoteResp-6]
	_ = x[MsgTypeAppendEntriesReq-7]
	_ = x[MsgTypeAcceptEntriesResp-8]
	_ = x[MsgTypeRefuseEntriesResp-9]
	_ = x[MsgTypeInstallSnapReq-10]
}

const (
	_MsgType_name_0 = "EventUknown"
	_MsgType_name_1 = "EventHeartbeatReqEventHeartbeatResponseEventClientReqEventClientResponseEventVoteReqEventVoteResponseEventAppendEntriesReqEventAcceptEntriesResponseEventRefuseEntriesResponseEventInstallSnapshot"
)

var (
	_MsgType_index_1 = [...]uint8{0, 17, 39, 53, 72, 84, 101, 122, 148, 174, 194}
)

func (i MsgType) String() string {
	switch {
	case i == -1:
		return _MsgType_name_0
	case 1 <= i && i <= 10:
		i -= 1
		return _MsgType_name_1[_MsgType_index_1[i]:_MsgType_index_1[i+1]]
	default:
		return "MsgType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
