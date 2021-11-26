package raft

import "fmt"

//RaftNode raft rsm node
type RaftNode struct {
	id         uint64
	term       uint64
	instC      chan Instruction //
	msgC       chan Message
	log        *RaftLog
	queuedReqs []queuedEvent
	proxyReqs  map[ReqId]Address
	peers      []string
	role       RaftRole
}

//NewRaftNode allocate a new RaftNode struct in heap and init
func NewRaftNode(id uint64,
	role RoleType,
	logStore LogStore,
	instc chan Instruction,
	msgC chan Message,
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
		node.role = NewFollower(node)
	case RoleLeader:
		node.role = NewLeader(node)
	}

	return node
}

func (node *RaftNode) String() string {
	return fmt.Sprintf("Node(%v) term(%v) role(%v)", node.id, node.term, node.RoleType())
}

//RoleType get raftnode role type
func (node *RaftNode) RoleType() RoleType {
	return node.role.Type()
}

func (node *RaftNode) becomeRole(roleType RoleType) {
	if node.RoleType() == roleType {
		return
	}

	logger.Debug("node(%v) become role %v\n", node, roleType)

	switch roleType {
	case RoleCandidate:
		node.role = NewCandidate(node)
	case RoleFollower:
		node.role = NewFollower(node)
	case RoleLeader:
		node.role = NewLeader(node)
	}
}

func (node *RaftNode) saveTermVoteMeta(term uint64, voteFor uint64) {}

func (node *RaftNode) becomeCandidate() {
	if node.RoleType() != RoleFollower {
		logger.Error("Node(%v) becomeCandidate \n", node)
		return
	}
	node.term++
	node.log.SaveTerm(node.term, "")

	lastIndex, lastTerm := node.log.LastIndexTerm()
	node.send(&AddrPeers{}, &EventSolicitVoteReq{lastIndex, lastTerm})

	node.becomeRole(RoleCandidate)
}

func (node *RaftNode) becomeFollower(term uint64, leader string) *RaftNode {
	if term < node.term {
		logger.Error("Node(%v) becomeFollower err, term: %v, leader: %v\n", node, term, leader)
		return node
	}
	node.term = term
	node.log.SaveTerm(node.term, leader)
	node.becomeRole(RoleFollower)
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
	committedIndex, committedTerm := node.log.CommittedIndexTerm()
	heartbeatEvent := &EventHeartbeatReq{
		commitIndex: committedIndex,
		commitTerm:  committedTerm,
	}
	node.send(AddressPeers, heartbeatEvent)

	node.appendAndCastCommand(NOOPCommand)
	node.abortProxyReqs()

	return node
}

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
		if queueReq.EventType() != EventTypeClientReq {
			logger.Warn("forward req(%v) \n", queueReq)
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

		msg := Message{
			from:  proxyFrom,
			to:    leader,
			term:  0,
			event: queueReq.event,
		}
		node.msgC <- msg
	}
}

// send a appendEntry
func (node *RaftNode) replicate(peer string, entry ...Entry) {
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
	msg := Message{
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
		logger.Warn("Node(%v) recieve invalid msg: %v", node, msg)
		return
	}

	// msg from peer which term > self
	if msg.term > node.term && msg.from.Type() == AddrTypePeer {
		from := msg.from.(*AddrPeer)
		node.becomeFollower(msg.term, from.peer)
	}

	node.role.Step(msg)
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
