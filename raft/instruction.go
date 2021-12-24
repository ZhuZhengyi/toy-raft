//
package raft

import "github.com/google/uuid"

//go:generate stringer -type=InstType

type InstType int16 // Instruction type

const (
	InstTypeUnkown InstType = -1   // Unknown Instruction
	InstTypeAbort  InstType = iota // Abort
	InstTypeApply                  //
	InstTypeNotify                 //
	InstTypeQuery                  //
	InstTypeStatus                 //
	InstTypeVote                   //
)

//Instruction: represent some data which outer client send to rsm
// to drive rsm to run
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
	id    ReqId
	index uint64
	addr  Address
}

type InstQuery struct {
	id      uuid.UUID
	term    uint64
	index   uint64
	quorum  uint64
	addr    Address
	command []byte
}

type InstStatus struct {
	id      uint64
	addr    Address
	command []byte
}

//
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
