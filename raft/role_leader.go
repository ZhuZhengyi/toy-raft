package raft

var (
	_ RaftNode = (*Leader)(nil)
)

type Leader struct {
	*raftNode
	heartbeatTicks   uint64
	heartbeatTimeOut uint64
	peerNextIndex    map[uint64]uint64
	peerLastIndex    map[uint64]uint64
}

func NewLeader(r *raftNode) *Leader {
	lastIndex, _ := r.log.LastIndexTerm()
	l := &Leader{
		raftNode:         r,
		heartbeatTicks:   0,
		heartbeatTimeOut: 1,
		peerNextIndex:    make(map[uint64]uint64),
		peerLastIndex:    make(map[uint64]uint64),
	}

	for _, peer := range r.peers {
		l.peerNextIndex[peer] = lastIndex + 1
		l.peerLastIndex[peer] = 0
	}

	return l
}

func (l *Leader) RoleType() RoleType {
	return RoleLeader
}

func (l *Leader) Step(msg *Message) {
}

func (l *Leader) Tick() {
	if len(l.peers) > 0 {
		l.heartbeatTicks += 1
		if l.heartbeatTicks >= l.heartbeatTimeOut {
			l.heartbeatTicks = 0

			commitIndex, commitTerm := l.log.CommittedIndexTerm()
			l.send(AddressPeers, &EventHeartbeatReq{commitIndex, commitTerm})
		}
	}
}
