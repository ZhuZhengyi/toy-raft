package raft

import "fmt"

//RaftNode raft rsm node
type RaftNode struct {
	id         uint64
	term       uint64
	instC      chan Instruction //
	msgC       chan *Message
	log        *RaftLog
	queuedReqs []queuedEvent
	proxyReqs  map[ReqId]Address
	peers      []uint64
	role       RaftRole
}

//NewRaftNode allocate a new RaftNode struct in heap and init
func NewRaftNode(id uint64,
	role RoleType,
	logStore LogStore,
	instc chan Instruction,
	msgC chan *Message,
) *RaftNode {
	raftLog := NewRaftLog(logStore)
	term, _ := raftLog.LoadTerm()
	node := &RaftNode{
		id:         id,
		term:       term,
		instC:      instc,
		msgC:       msgC,
		queuedReqs: make([]queuedEvent, QUEUED_REQ_BATCH_SIZE),
		proxyReqs:  make(map[ReqId]Address),
		log:        raftLog,
	}

	switch role {
	case RoleCandidate:
		node.role = NewCandidate(node)
	case RoleFollower:
		node.role = NewFollower(node, 0, 0)
	case RoleLeader:
		node.role = NewLeader(node)
	}

	return node
}

func (node *RaftNode) String() string {
	return fmt.Sprintf("{id:%v,term:%v,role:%v}", node.id, node.term, node.RoleType())
}

//RoleType get raftnode role type
func (node *RaftNode) RoleType() RoleType {
	return node.role.Type()
}

func (node *RaftNode) becomeRole(roleType RoleType) {
	if node.RoleType() == roleType {
		return
	}

	logger.Info("node:%v change role: %v -> %v", node, node.role.Type(), roleType)

	switch roleType {
	case RoleCandidate:
		node.role = NewCandidate(node)
	case RoleFollower:
		node.role = NewFollower(node, 0, 0)
	case RoleLeader:
		node.role = NewLeader(node)
	}
}

func (node *RaftNode) saveTermVoteMeta(term uint64, voteFor uint64) {}

func (node *RaftNode) becomeCandidate() {
	switch node.role.(type) {
	case *Leader:
		logger.Warn("node:%v can't becomeCandidate \n", node)
		return
	}

	node.becomeRole(RoleCandidate)

	node.term++
	node.log.SaveTerm(node.term, 0)

	node.send(node.AddrPeers(), &EventSolicitVoteReq{node.log.lastIndex, node.log.lastTerm})
}

func (node *RaftNode) AddrPeers() *AddrPeers {
	return &AddrPeers{peers: node.peers}
}

func (node *RaftNode) becomeFollower(term uint64, leader uint64) *RaftNode {
	votedFor := uint64(0)
	if term < node.term {
		logger.Error("node:%v becomeFollower err, term: %v, leader: %v\n", node, term, leader)
		return node
	}

	logger.Debug("node:%v becomeFollower term: %v, leader: %v\n", node, term, leader)

	switch r := node.role.(type) {
	case *Follower:
		if term > node.term {
			logger.Info("node:%v discover a new term:%v, following leader %v", node, term, leader)
			node.term = term
			node.log.SaveTerm(node.term, 0)
		} else {
			logger.Info("node:%v discover a new leader:%v, following", node, leader)
		}
		votedFor = r.votedFor
		node.abortProxyReqs()
		node.forwardToLeaderQueued(&AddrPeer{leader})
		node.role = NewFollower(node, leader, votedFor)
	case *Candidate:
		node.term = term
		node.log.SaveTerm(term, 0)
		node.abortProxyReqs()
		node.forwardToLeaderQueued(&AddrPeer{leader})
		node.role = NewFollower(node, leader, 0)
	case *Leader:
		node.term = term
		node.log.SaveTerm(term, 0)
		node.instC <- &InstAbort{}
		node.role = NewFollower(node, leader, 0)
	default:
	}

	return node
}

func (node *RaftNode) becomeLeader() *RaftNode {
	switch node.role.(type) {
	case *Follower:
		logger.Warn("node:%v cannot allow change to leader", node)
		return node
	}

	logger.Detail("node:%v will change %v -> Leader", node, node.role.Type())
	node.becomeRole(RoleLeader)

	heartbeatEvent := &EventHeartbeatReq{
		commitIndex: node.log.commitIndex,
		commitTerm:  node.log.commitTerm,
	}
	node.send(&AddrPeers{peers: node.peers}, heartbeatEvent)

	node.appendAndCastCommand(NOOPCommand)
	node.abortProxyReqs()

	logger.Info("node:%v Leader:[term:%v,leader:%v]", node, node.term, node.id)
	return node
}

//append command
func (node *RaftNode) appendAndCastCommand(command []byte) {
	entry := node.log.Append(node.term, command)

	for _, p := range node.peers {
		node.replicate(p, *entry)
	}
}

func (node *RaftNode) quorum() uint64 {
	return (uint64(len(node.peers))+1)/2 + 1
}

func (node *RaftNode) abortProxyReqs() {
	for id, addr := range node.proxyReqs {
		node.send(addr, &EventClientResp{id: id})
	}
}

// forward to leader
func (node *RaftNode) forwardToLeaderQueued(leader Address) {
	switch node.role.(type) {
	case *Leader:
		logger.Warn("node:%v cannot forwardToLeaderQueued", node)
		return
	}
	for _, queueReq := range node.queuedReqs {
		switch event := queueReq.event.(type) {
		case *EventClientReq:
			from := queueReq.from
			node.proxyReqs[event.id] = from
			proxyFrom := from
			switch from.(type) {
			case *AddrClient:
				proxyFrom = AddressLocal
			}
			msg := &Message{
				from:  proxyFrom,
				to:    leader,
				term:  0,
				event: event,
			}
			node.msgC <- msg
			continue
		}
	}
}

// send a appendEntry
func (node *RaftNode) replicate(peer uint64, entry ...Entry) {
	appendEntriesEvent := &EventAppendEntriesReq{
		baseIndex: node.log.lastIndex,
		baseTerm:  node.term,
		entries:   make([]Entry, 0),
	}
	appendEntriesEvent.entries = append(appendEntriesEvent.entries, entry...)
	node.send(&AddrPeer{peer: peer}, appendEntriesEvent)
}

// send event to
func (node *RaftNode) send(to Address, event MsgEvent) {
	msg := &Message{
		from:  AddressLocal,
		to:    to,
		term:  node.term,
		event: event,
	}
	node.msgC <- msg
	logger.Detail("node:%v send msg:%v", node, msg)
}

//Step step rsm by msg
func (node *RaftNode) Step(msg *Message) {
	if !node.validateMsg(msg) {
		logger.Debug("node:%v,step_invalid_msg:%v", node, msg)
		return
	}

	logger.Detail("node:%v,step msg:%v", node, msg)

	// msg from peer which term > self
	switch from := msg.from.(type) {
	case *AddrPeer:
		if msg.term > node.term || node.isFollowerNoLeader() {
			logger.Debug("node:%v -> follower, msg:%v term is bigger", node, msg)
			node.becomeFollower(msg.term, from.peer)
		}
	default:
	}

	node.role.Step(msg)
}

func (node *RaftNode) isFollowerNoLeader() bool {
	switch role := node.role.(type) {
	case *Follower:
		return role.leader == 0
	}
	return false
}

//Tick tick rsm
func (node *RaftNode) Tick() {
	node.role.Tick()
}

func (node *RaftNode) validateMsg(msg *Message) bool {
	switch msg.from.(type) {
	case *AddrPeers, *AddrLocal:
		logger.Error("node:%v msg:%v from addr invalid", node, msg)
		return false
	case *AddrClient:
		switch msg.event.(type) {
		case *EventClientReq:
		default:
			logger.Error("node:%v msg:%v from client invalid event", node, msg)
			return false
		}
	}

	if msg.term < node.term {
		switch msg.event.(type) {
		case *EventClientReq, *EventClientResp:
		default:
			logger.Debug("node:%v msg:%v invalid event", node, msg)
			return false
		}
	}

	// check msg.to
	switch to := msg.to.(type) {
	case *AddrPeer:
		if to.peer == node.id {
			return true
		} else {
			logger.Error("node:%v,msg:%v to peer id is invalid", node, msg)
			return false
		}
	case *AddrLocal, *AddrPeers:
		return true
	case *AddrClient:
		return false
	}

	return false
}
