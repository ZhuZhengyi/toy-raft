//
package raft

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
