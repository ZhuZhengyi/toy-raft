// request.go
package raft

//go:generate stringer -type=ReqType,RespType  -linecomment

type (
	RespType int32
	ReqType  int32
)

const (
	ReqTypeQuery  ReqType = iota // ReqQuery
	ReqTypeMutate                // ReqMutate
	ReqTypeStatus                // ReqStatus

	RespTypeStatus RespType = iota // RespStatus
	RespTypeState                  // RespState
)

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
