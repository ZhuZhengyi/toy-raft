// instruct_state.go
package raft

// 指令状态机
type InstStateMachine interface {
	AppliedIndex() uint64
	Apply(index uint64, command []byte)
	Query(command []byte)
}

type DummyInstStateMachine struct {
}

func (sm *DummyInstStateMachine) AppliedIndex() (index uint64) {
	return
}

func (sm *DummyInstStateMachine) Apply(index uint64, command []byte) {
}

func (sm *DummyInstStateMachine) Query(command []byte) {
}
