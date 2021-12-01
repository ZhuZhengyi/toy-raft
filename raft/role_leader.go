package raft

type Leader struct {
	*RaftNode
	heartbeatTicks   uint64
	heartbeatTimeOut uint64
	peerNextIndex    map[uint64]uint64
	peerLastIndex    map[uint64]uint64
}

func NewLeader(node *RaftNode) *Leader {
	lastIndex, _ := node.log.LastIndexTerm()
	l := &Leader{
		RaftNode:         node,
		heartbeatTicks:   0,
		heartbeatTimeOut: 1,
		peerNextIndex:    make(map[uint64]uint64),
		peerLastIndex:    make(map[uint64]uint64),
	}

	for _, peer := range node.peers {
		l.peerNextIndex[peer] = lastIndex + 1
		l.peerLastIndex[peer] = 0
	}

	return l
}

var (
	_ RaftRole = (*Leader)(nil)
)

//Type leader type
func (l *Leader) Type() RoleType {
	return RoleLeader
}

//Step step rsm by msg
func (l *Leader) Step(msg *Message) {
	switch msg.MsgType() {
	case MsgTypeHeartbeatResp:
		//TODO:
	case MsgTypeAcceptEntriesResp:
		//TODO:
	case MsgTypeRefuseEntriesResp:
		//TODO:
	case MsgTypeClientReq:
		//TODO:
	case MsgTypeClientResp:
		//TODO:
	case MsgTypeVoteReq, MsgTypeVoteResp:
		{
		}
	default:
		logger.Warn("role(%v) step unexpected msg(%v)\n", l, msg)
		//TODO:
	}
}

//Tick tick leader
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
