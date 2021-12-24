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
		heartbeatTimeOut: HB_LEADER_TIMEOUT,
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
	logger.Detail("leader:%v,step msg:%v", l, msg)
	switch event := msg.event.(type) {
	case *EventHeartbeatResp:
		switch from := msg.from.(type) {
		case *AddrPeer:
			instruction := &InstVote{term: msg.term, index: event.commitIndex, addr: msg.from}
			l.instC <- instruction
			if !event.hasCommitted {
				l.replicate(from.peer)
			}
		}
	case *EventAcceptEntriesResp:
		switch from := msg.from.(type) {
		case *AddrPeer:
			l.peerLastIndex[from.peer] = event.lastIndex
			l.peerNextIndex[from.peer] = event.lastIndex + 1
		}
		l.commit()
	case *EventRefuseEntriesResp:
		switch from := msg.from.(type) {
		case *AddrPeer:
			if i := l.peerNextIndex[from.peer]; i > 1 {
				l.peerLastIndex[from.peer] = i - 1
			}
			l.replicate(from.peer)
		}
	case *EventClientReq:
		switch req := event.request.(type) {
		case *ReqQuery:
			instQuery := &InstQuery{id: event.id,
				addr:    msg.from,
				command: req.command,
				term:    l.term,
				index:   l.log.CommittedIndex(),
				quorum:  l.quorum(),
			}
			l.instC <- instQuery
			instVote := &InstVote{term: l.term}
			l.instC <- instVote
			if len(l.peers) > 0 {
				index, term := l.log.CommittedIndexTerm()
				l.send(&AddrPeers{peers: l.peers}, &EventHeartbeatReq{
					commitIndex: index,
					commitTerm:  term,
				})
			}
		case *ReqMutate:
			index := l.append(req.mutate)
			l.instC <- &InstNotify{id: event.id, addr: msg.from, index: index}
			if len(l.peers) == 0 {
				l.commit()
			}
		case *ReqStatus:
		default:
		}
	case *EventClientResp:
		//TODO:
		l.send(&AddrClient{}, event)
	case *EventSolicitVoteReq, *EventGrantVoteResp:
		logger.Warn("leader:%v ignore msg(%v)\n", l, msg)
	case *EventHeartbeatReq, *EventAppendEntriesReq:
		logger.Warn("leader:%v step unexpected msg(%v)\n", l, msg)
	default:
		logger.Warn("leader:%v step unexpected msg(%v)\n", l, msg)
	}
}

//Tick tick leader
func (l *Leader) Tick() {
	logger.Detail("leader:%v tick", l)
	if len(l.peers) > 0 {
		l.heartbeatTicks += 1
		if l.heartbeatTicks >= l.heartbeatTimeOut {
			logger.Detail("leader:%v hbtick timeout", l)
			l.heartbeatTicks = 0
			commitIndex, commitTerm := l.log.CommittedIndexTerm()
			l.send(AddressPeers, &EventHeartbeatReq{commitIndex, commitTerm})
		}
	}
}

/// replicate logs to peer
func (l *Leader) replicate(peer uint64) {
	//TODO: replicate logs to peer
	logger.Detail("leader:%v replicate log to peer:%v", l, peer)

}

/// Commits pending log entries
func (l *Leader) commit() {
	logger.Detail("leader:%v commit pending entries", l)
	//TODO
}

// append entry to log and replicate to peers
func (l *Leader) append(command []byte) uint64 {
	logger.Detail("leader:%v append command to log", l)
	entry := l.log.Append(l.term, command)
	for _, peer := range l.peers {
		l.replicate(peer)
	}

	return entry.index
}
