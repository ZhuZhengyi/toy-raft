//
package raft

//go:generate stringer -type=InstType

type InstType int16

const (
	InstTypeUnkown InstType = -1
	InstTypeAbort  InstType = iota
	InstTypeApply
	InstTypeNotify
	InstTypeQuery
	InstTypeStatus
	InstTypeVote
)

type Instruction interface {
	Type() InstType
}

var (
	_ Instruction = (*InstAbort)(nil)
	_ Instruction = (*InstApply)(nil)
	_ Instruction = (*InstNotify)(nil)
	_ Instruction = (*InstQuery)(nil)
	_ Instruction = (*InstStatus)(nil)
	_ Instruction = (*InstVote)(nil)
)

type InstAbort struct {
}

type InstApply struct {
	entry *Entry
}

type InstNotify struct {
	id    uint64
	index uint64
	addr  Address
}

type InstQuery struct {
	id      uint64
	term    uint64
	quorum  uint64
	addr    Address
	command []byte
}

type InstStatus struct {
	id      uint64
	addr    Address
	command []byte
}

type InstVote struct {
	term  uint64
	index uint64
	addr  Address
}

func (i *InstAbort) Type() InstType {
	return InstTypeAbort
}

func (i *InstStatus) Type() InstType {
	return InstTypeStatus
}

func (i *InstVote) Type() InstType {
	return InstTypeVote
}

func (i *InstQuery) Type() InstType {
	return InstTypeQuery
}

func (i *InstApply) Type() InstType {
	return InstTypeApply
}

func (i *InstNotify) Type() InstType {
	return InstTypeNotify
}
