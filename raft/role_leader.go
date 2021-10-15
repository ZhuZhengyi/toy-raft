package raft

type leader struct {
	raftNode
	heartbeatTicks uint64
	peerNextIndex  map[string]uint64
	peerLastIndex  map[string]uint64
}

var (
	_ RaftRole = (*leader)(nil)
)

func NewLeader(node *raftNode) *leader {
	l := &leader{
		raftNode:       *node,
		heartbeatTicks: 0,
		peerNextIndex:  make(map[string]uint64),
		peerLastIndex:  make(map[string]uint64),
	}

	return l
}

func (l *leader) RoleType() RoleType {
	return RoleLeader
}

func (l *leader) Step(msg *Message) {
}

func (l *leader) Tick() {
	l.heartbeatTicks += 1
}
