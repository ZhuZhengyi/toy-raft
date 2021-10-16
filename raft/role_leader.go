package raft

type leader struct {
	*raftNode
	heartbeatTicks uint64
	peerNextIndex  map[uint64]uint64
	peerLastIndex  map[uint64]uint64
}

var (
	_ RaftRole = (*leader)(nil)
)

func NewLeader(node *raftNode) *leader {
	lastIndex, _ := node.log.LastIndexTerm()
	l := &leader{
		raftNode:       node,
		heartbeatTicks: 0,
		peerNextIndex:  make(map[uint64]uint64),
		peerLastIndex:  make(map[uint64]uint64),
	}

	for _, peer := range node.peers {
		l.peerNextIndex[peer] = lastIndex + 1
		l.peerLastIndex[peer] = 0
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
