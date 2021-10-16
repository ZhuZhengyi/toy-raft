// request.go
package raft

type (
	RespType int32
	ReqType  int32
)

const (
	ReqTypeUnknown ReqType = -1
	ReqTypeQuery   ReqType = iota
	ReqTypeMutate
	ReqTypeStatus

	RespTypeUnknown RespType = -1
	RespTypeStatus  RespType = iota
	RespTypeState
)

func (r ReqType) String() string {
	switch r {
	case ReqTypeQuery:
		return "ReqQuery"
	case ReqTypeMutate:
		return "ReqMutate"
	case ReqTypeStatus:
		return "ReqStatus"
	default:
		return "ReqUnkown"
	}
}

func (r RespType) String() string {
	switch r {
	case RespTypeStatus:
		return "RespStatus"
	case RespTypeState:
		return "RespState"
	default:
		return "RespUnkown"
	}
}

type Request interface {
	Type() ReqType
}

type Response interface {
	Type() RespType
}

type Session interface {
	Send(Response)
}

var (
	_ Request = (*ReqQuery)(nil)
	_ Request = (*ReqMutate)(nil)
	_ Request = (*ReqStatus)(nil)

	_ Response = (*RespState)(nil)
	_ Response = (*RespStatus)(nil)
)

type ReqQuery struct {
	query []byte
}

type ReqMutate struct {
	mutate []byte
}

type ReqStatus struct {
}

type RespStatus struct {
}

type RespState struct {
	state []byte
}

func (r *ReqQuery) Type() ReqType {
	return ReqTypeQuery
}

func (r *ReqMutate) Type() ReqType {
	return ReqTypeMutate
}

func (r *ReqStatus) Type() ReqType {
	return ReqTypeStatus
}

func (r *RespState) Type() RespType {
	return RespTypeState
}

func (r *RespStatus) Type() RespType {
	return RespTypeStatus
}
