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

	logger.Info("node:%v become role %v\n", node, roleType)

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
	if node.RoleType() == RoleLeader {
		logger.Error("Node(%v) can't becomeCandidate \n", node)
		return
	}

	node.term++
	node.log.SaveTerm(node.term, 0)

	lastIndex, lastTerm := node.log.LastIndexTerm()
	node.becomeRole(RoleCandidate)

	node.send(node.AddrPeers(), &EventSolicitVoteReq{lastIndex, lastTerm})
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
	if term > node.term {
		logger.Info("node:%v discover a new term:%v, following leader %v", node, term, leader)
		node.term = term
		node.log.SaveTerm(node.term, 0)
	} else {
		logger.Info("node:%v discover a new leader:%v, following", node, leader)
		//votedFor = node.role
	}

	switch r := node.role.(type) {
	case *Follower:
		votedFor = r.votedFor
		node.role = NewFollower(node, leader, votedFor)
	case *Candidate:

	case *Leader:
	default:

	}

	node.abortProxyReqs()
	node.forwardToLeaderQueued(&AddrPeer{leader})

	return node
}

func (node *RaftNode) becomeLeader() *RaftNode {
	if node.RoleType() != RoleCandidate {
		logger.Warn("becomeLeader Warn  ")
		return node
	}

	node.becomeRole(RoleLeader)

	logger.Info("node:%v become leader", node)

	committedIndex, committedTerm := node.log.CommittedIndexTerm()
	heartbeatEvent := &EventHeartbeatReq{
		commitIndex: committedIndex,
		commitTerm:  committedTerm,
	}
	node.send(&AddrPeers{peers: node.peers}, heartbeatEvent)

	node.appendAndCastCommand(NOOPCommand)
	node.abortProxyReqs()

	for _, queuedReq := range node.queuedReqs {
		node.Step(&Message{
			from:  queuedReq.from,
			to:    &AddrLocal{},
			term:  0,
			event: queuedReq.event,
		})
	}

	logger.Info("node:%v become leader", node)
	return node
}

//append command
func (node *RaftNode) appendAndCastCommand(command []byte) {
	entry := node.log.Append(node.term, command)

	for _, p := range node.peers {
		node.replicate(p, entry)
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
	if node.RoleType() == RoleLeader {
		return
	}
	for _, queueReq := range node.queuedReqs {
		if queueReq.MsgType() != MsgTypeClientReq {
			//logger.Warn("forward req(%v) \n", queueReq)
			continue
		}
		originEvent := queueReq.event.(*EventClientReq)
		from := queueReq.from

		// record origin req
		node.proxyReqs[originEvent.id] = from

		// forward
		proxyFrom := from
		if from.Type() == AddrTypeClient {
			proxyFrom = AddressLocal
		}

		msg := &Message{
			from:  proxyFrom,
			to:    leader,
			term:  0,
			event: queueReq.event,
		}
		node.msgC <- msg
	}
}

// send a appendEntry
func (node *RaftNode) replicate(peer uint64, entry ...Entry) {
	index, _ := node.log.LastIndexTerm()
	appendEntriesEvent := &EventAppendEntriesReq{
		baseIndex: index,
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
			logger.Detail("node:%v,msg:%v term is bigger", node, msg)
			node.becomeFollower(msg.term, from.peer)
		}
	default:
	}

	node.role.Step(msg)
}

func (node *RaftNode) isFollowerNoLeader() bool {
	if node.RoleType() == RoleFollower {
		if n, ok := node.role.(*Follower); ok {
			return n.leader == 0
		}
	}
	return false
}

//Tick tick rsm
func (node *RaftNode) Tick() {
	node.role.Tick()
}

func (node *RaftNode) validateMsg(msg *Message) bool {
	if msg.term < node.term {
		return false
	}

	//TODO:

	return true
}
