//
package raft

type InstType int16

const (
	InstTypeUnkown InstType = -1

	InstTypeAbort InstType = iota
	InstTypeApply
	InstTypeNotify
	InstTypeQuery
	InstTypeStatus
	InstTypeVote

	InstNameUnkown = "InstUnkown"
	InstNameAbort  = "InstAbort"
	InstNameApply  = "InstApply"
	InstNameNotify = "InstNotify"
	InstNameQuery  = "InstQuery"
	InstNameStatus = "InstStatus"
	InstNameVote   = "InstVote"
)

type Instruction interface {
	Type() InstType
	String() string
}

type instruction struct {
}

func (i *instruction) Type() InstType {
	return InstTypeUnkown
}

func (i *instruction) String() string {
	return InstNameUnkown
}

type InstAbort struct {
	instruction
}

type InstApply struct {
	instruction
	entry Entry
}

type InstNotify struct {
	instruction
	id    uint64
	index uint64
	addr  Address
}

type InstQuery struct {
	instruction
	id      uint64
	term    uint64
	quorum  uint64
	addr    Address
	command []byte
}

type InstStatus struct {
	instruction
	id      uint64
	addr    Address
	command []byte
}

type InstVote struct {
	instruction
	term  uint64
	index uint64
	addr  Address
}
