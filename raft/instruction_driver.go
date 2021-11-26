//
package raft

import "context"

type InstDriver struct {
	stopC chan struct{}
	instC chan Instruction
	msgC  chan Message
	//notify map[uint64]struct{addr Address, data []byte}
	sm InstStateMachine
}

//NewInstDriver allocate a new InstDriver from heap
func NewInstDriver(instC chan Instruction, msgC chan Message, sm InstStateMachine) *InstDriver {
	return &InstDriver{
		stopC: make(chan struct{}),
		instC: instC,
		msgC:  msgC,
		sm:    sm,
	}
}

func (d *InstDriver) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		case i := <-d.instC:
			d.execute(i)
		}
	}
}

//
func (d *InstDriver) stop() {
	d.stopC <- struct{}{}
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

// execute instruction
func (d *InstDriver) execute(inst Instruction) {
	switch inst.(type) {
	case *InstAbort:
		d.queryAbort()
	case *InstVote:
		d.queryVote()
	case *InstApply:
		applyInst := inst.(*InstApply)
		entry := applyInst.entry
		d.sm.Apply(entry.index, entry.command)
	default:
		logger.Warn("execute unexpected inst(%v)\n", inst)
	}
}
