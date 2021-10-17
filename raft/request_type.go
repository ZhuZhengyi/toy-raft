// request.go
package raft

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
