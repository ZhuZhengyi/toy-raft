// request.go
package raft

type Session interface {
	Send(Response)
}

type session struct {
}

func (s *session) Send(Response) {
}

type Request interface {
	Type() ReqType
	Session() Session
}

type Response interface {
	Type() RespType
}

type (
	RespType int32
	ReqType  int32
)

const (
	RespTypeUnknown RespType = -1

	RespTypeStatus RespType = iota
	RespTypeState
)

const (
	ReqTypeUnknown ReqType = -1

	ReqTypeQuery ReqType = iota
	ReqTypeMutate
	ReqTypeStatus
)

var (
	_ Request = new(request)
)

type request struct {
	session *session
}

func (r *request) Session() Session {
	return r.session
}

func (r *request) Type() ReqType {
	switch Request(r).(type) {
	case *ReqQuery:
		return ReqTypeQuery
	case *ReqMutate:
		return ReqTypeMutate
	case *ReqStatus:
		return ReqTypeStatus
	}

	return ReqTypeUnknown
}

type ReqQuery struct {
	request
	query []byte
}

type ReqMutate struct {
	request
	mutate []byte
}

type ReqStatus struct {
	request
}

type response struct {
}

func (r *response) Type() RespType {
	switch Response(r).(type) {
	case *RespStatus:
		return RespTypeStatus
	case *RespState:
		return RespTypeState
	}

	return RespTypeUnknown
}

type RespStatus struct {
	response
}

type RespState struct {
	response
	state []byte
}
