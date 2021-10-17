package raft

import "fmt"

type RaftNode struct {
	id         uint64
	term       uint64
	instC      chan Instruction
	msgC       chan Message
	log        *RaftLog
	queuedReqs []queuedEvent
	proxyReqs  map[ReqId]Address
	peers      []uint64
	role       RaftRole
}

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

func (r *RaftNode) String() string {
	return fmt.Sprintf("Node(%v) term(%v)", r.id, r.term)
}

func (r *RaftNode) RoleType() RoleType {
	return r.role.Type()
}

func (r *RaftNode) becomeRole(roleType RoleType) {
	logger.Debug("node(%v) change role %v -> %v\n", r.id, r.role.Type(), roleType)

	if r.RoleType() == roleType {
		return
	}
	switch roleType {
	case RoleCandidate:
		r.role = NewCandidate(r)
	case RoleFollower:
		r.role = NewFollower(r)
	case RoleLeader:
		r.role = NewLeader(r)
	}
}

func (r *RaftNode) saveTermVoteMeta(term uint64, voteFor uint64) {}

func (r *RaftNode) becomeCandidate() {
	if r.RoleType() != RoleFollower {
		logger.Error("Node(%v) becomeCandidate \n", r)
		return
	}
	r.term += 1
	r.log.SaveTerm(r.term, 0)

	lastIndex, lastTerm := r.log.LastIndexTerm()
	r.send(&AddrPeers{}, &EventSolicitVoteReq{lastIndex, lastTerm})

	r.becomeRole(RoleCandidate)
}

func (r *RaftNode) becomeFollower(term, leader uint64) *RaftNode {
	if term < r.term {
		logger.Error("Node(%v) becomeFollower err, term: %v, leader: %v\n", r, term, leader)
		return r
	}
	r.term = term
	r.log.SaveTerm(r.term, leader)
	r.becomeRole(RoleFollower)
	r.abortProxyReqs()
	r.forwardToLeaderQueued(&AddrPeer{leader})
	return r
}

func (r *RaftNode) becomeLeader() *RaftNode {
	if r.RoleType() != RoleCandidate {
		logger.Warn("becomeLeader Warn  ")
		return r
	}

	r.becomeRole(RoleLeader)
	committedIndex, committedTerm := r.log.CommittedIndexTerm()
	heartbeatEvent := &EventHeartbeatReq{
		commitIndex: committedIndex,
		commitTerm:  committedTerm,
	}
	r.send(AddressPeers, heartbeatEvent)

	r.appendAndCastCommand(NOOPCommand)
	r.abortProxyReqs()

	return r
}

func (r *RaftNode) appendAndCastCommand(command []byte) {
	entry := r.log.Append(r.term, command)

	for _, p := range r.peers {
		r.replicate(p, entry)
	}
}

func (r *RaftNode) quorum() uint64 {
	return (uint64(len(r.peers))+1)/2 + 1
}

func (r *RaftNode) abortProxyReqs() {
	for id, addr := range r.proxyReqs {
		r.send(addr, &EventClientResp{id: id})
	}
}

func (r *RaftNode) forwardToLeaderQueued(leader Address) {
	if r.RoleType() == RoleLeader {
		return
	}
	for _, queuedEvent := range r.queuedReqs {
		if queuedEvent.event.Type() == EventTypeClientReq {
			originEvent := queuedEvent.event.(*EventClientReq)
			from := queuedEvent.from

			// record origin req
			r.proxyReqs[originEvent.id] = from

			// forward
			proxyFrom := from
			if from.Type() == AddrTypeClient {
				proxyFrom = new(AddrLocal)
			}

			msg := Message{
				from:  proxyFrom,
				to:    leader,
				term:  0,
				event: queuedEvent.event,
			}
			r.msgC <- msg
		}

	}
}

func (r *RaftNode) replicate(peer uint64, entry Entry) {
	appendEntriesEvent := &EventAppendEntriesReq{}
	r.send(&AddrPeer{peer: peer}, appendEntriesEvent)
}

func (r *RaftNode) send(to Address, event MsgEvent) {
	msg := Message{
		from:  new(AddrLocal),
		to:    to,
		term:  r.term,
		event: event,
	}
	r.msgC <- msg
}

func (r *RaftNode) Step(msg *Message) {
	if !r.validateMsg(msg) {
		return
	}

	// msg from peer which term > self
	if msg.term > r.term && msg.from.Type() == AddrTypePeer {
		from := msg.from.(*AddrPeer)
		r.becomeFollower(msg.term, from.peer)
	}

	r.role.Step(msg)
}

func (r *RaftNode) Tick() {
	r.role.Tick()
}

func (r *RaftNode) validateMsg(msg *Message) bool {
	if msg.term < r.term {
		return false
	}

	//TODO:

	return true
}
