// instruct_state.go
package raft

// 指令状态机
type InstStateMachine interface {
	AppliedIndex() uint64
	Apply(index uint64, command []byte)
	Query(command []byte)
}

type InstructionType uint8

const (
	InstTypeAbort InstructionType = iota
	InstTypeApply
	InstTypeNotify
	InstTypeQuery
	InstTypeStatus
	InstTypeVote
)

type Instruction interface {
	Type() InstructionType
}

type InstAbort struct {
}

func (i *InstAbort) Type() InstructionType {
	return InstTypeAbort
}

type InstApply struct {
	entry Entry
}

func (i *InstApply) Type() InstructionType {
	return InstTypeApply
}

type InstNotify struct {
	id    uint64
	index uint64
	addr  Address
}

func (i *InstNotify) Type() InstructionType {
	return InstTypeNotify
}

type InstQuery struct {
	id      uint64
	term    uint64
	quorum  uint64
	addr    Address
	command []byte
}

func (i *InstQuery) Type() InstructionType {
	return InstTypeQuery
}

type InstStatus struct {
	id      uint64
	addr    Address
	command []byte
}

func (i *InstStatus) Type() InstructionType {
	return InstTypeStatus
}

type InstVote struct {
	term  uint64
	index uint64
	addr  Address
}

func (i *InstVote) Type() InstructionType {
	return InstTypeVote
}

type InstDriver struct {
	stopC chan struct{}
	instC chan Instruction
	msgC  chan Message
	//notify map[uint64]struct{addr Address, data []byte}
	sm InstStateMachine
}

func NewInstDriver(instC chan Instruction, msgC chan Message, sm InstStateMachine) *InstDriver {
	return &InstDriver{
		stopC: make(chan struct{}),
		instC: instC,
		msgC:  msgC,
		sm:    sm,
	}
}

func (d *InstDriver) drive() {
	for {
		select {
		case <-d.stopC:
			break
		case i := <-d.instC:
			d.execute(i)
		}
	}
}

func (d *InstDriver) notifyAbort() {
}

func (d *InstDriver) notifyApplied(index uint64, result []byte) {
}

func (d *InstDriver) queryAbort() {
}
func (d *InstDriver) queryVote() {
}
func (d *InstDriver) queryExecute() {
}

func (d *InstDriver) execute(inst Instruction) {
	switch inst.(type) {
	case *InstAbort:
		d.queryAbort()
	case *InstVote:
		d.queryVote()

	}
}

type DummyInstStateMachine struct {
}

func NewDummyInstStateMachine() *DummyInstStateMachine {
	return new(DummyInstStateMachine)
}

func (sm *DummyInstStateMachine) AppliedIndex() (index uint64) {
	return
}

func (sm *DummyInstStateMachine) Apply(index uint64, command []byte) {
}

func (sm *DummyInstStateMachine) Query(command []byte) {
}
