package raft

type queuedEvent struct {
	from  Address
	event MsgEvent
}

type RaftNode interface {
	RoleType() RoleType
	Step(msg *Message)
	Tick()
}

var (
	_ RaftNode = (*raftNode)(nil)
)

type raftNode struct {
	id         uint64
	term       uint64
	instC      chan Instruction
	msgC       chan Message
	log        *RaftLog
	queuedReqs []queuedEvent
	proxyReqs  map[ReqId]Address
	peers      []uint64
	role       RaftNode
}

func NewRaftNode(id uint64,
	role RoleType,
	logStore LogStore,
	instc chan Instruction,
	msgC chan Message,
) *raftNode {
	raftLog := NewRaftLog(logStore)
	term, _ := raftLog.LoadTerm()
	node := &raftNode{
		id:         id,
		term:       term,
		instC:      instc,
		msgC:       msgC,
		queuedReqs: make([]queuedEvent, 4096),
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

func (r *raftNode) RoleType() RoleType {
	return r.role.RoleType()
}

func (r *raftNode) becomeRole(roleType RoleType) {
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

func (r *raftNode) saveTermVoteMeta(term uint64, voteFor uint64) {}

func (r *raftNode) becomeCandidate() {
	if r.RoleType() != RoleFollower {
		return
	}
	r.term += 1
	r.log.SaveTerm(r.term, 0)

	lastIndex, lastTerm := r.log.LastIndexTerm()
	r.send(&AddrPeers{}, &EventSolicitVoteReq{lastIndex, lastTerm})

	r.becomeRole(RoleCandidate)
}

func (r *raftNode) becomeFollower(term, leader uint64) *raftNode {
	r.term = term
	r.log.SaveTerm(r.term, leader)
	r.becomeRole(RoleFollower)
	r.abortProxyReqs()
	r.forwardToLeaderQueued(&AddrPeer{leader})
	return r
}

func (r *raftNode) becomeLeader() *raftNode {
	if r.RoleType() != RoleCandidate {
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

func (r *raftNode) appendAndCastCommand(command []byte) {
	entry := r.log.Append(r.term, command)

	for _, p := range r.peers {
		r.replicate(p, entry)
	}
}

func (r *raftNode) quorum() uint64 {
	return (uint64(len(r.peers))+1)/2 + 1
}

func (r *raftNode) abortProxyReqs() {
	for id, addr := range r.proxyReqs {
		r.send(addr, &EventClientResp{id: id})
	}
}

func (r *raftNode) forwardToLeaderQueued(leader Address) {
	if r.RoleType() == RoleLeader {
		return
	}
	for _, queuedEvent := range r.queuedReqs {
		if queuedEvent.event.Type() == MsgTypeClientReq {
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

func (r *raftNode) replicate(peer uint64, entry Entry) {
	appendEntriesEvent := &EventAppendEntriesReq{}
	r.send(&AddrPeer{peer: peer}, appendEntriesEvent)
}

func (r *raftNode) send(to Address, event MsgEvent) {
	msg := Message{
		from:  new(AddrLocal),
		to:    to,
		term:  r.term,
		event: event,
	}
	r.msgC <- msg
}

func (r *raftNode) Step(msg *Message) {
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

func (r *raftNode) Tick() {
	r.role.Tick()
}

func (r *raftNode) validateMsg(msg *Message) bool {
	if msg.term < r.term {
		return false
	}

	return true
}
