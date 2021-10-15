package raft

type RoleType uint8

const (
	RoleFollower RoleType = iota
	RoleCandidate
	RoleLeader

	RoleNameFollower  = "Follower"
	RoleNameCandidate = "Candidate"
	RoleNameLeader    = "Leader"
	RoleNameUnkown    = "Unkown"
)

func (r RoleType) String() string {
	switch r {
	case RoleFollower:
		return RoleNameFollower
	case RoleCandidate:
		return RoleNameCandidate
	case RoleLeader:
		return RoleNameLeader
	default:
		return RoleNameUnkown
	}
}

type RaftRole interface {
	RoleType() RoleType
	Step(msg *Message)
	Tick()
}

type raftNode struct {
	id        uint64
	term      uint64
	instC     chan Instruction
	msgC      chan Message
	raftLog   *RaftLog
	requests  []Request
	proxyReqs map[uint64]Request
	role      RaftRole
}

var (
	_ RaftRole = (*raftNode)(nil)
)

func newRaftNode(id uint64, role RoleType, logStore LogStore,
	instc chan Instruction,
	msgC chan Message,
) *raftNode {
	raftLog := NewRaftLog(logStore)
	term, _ := raftLog.LoadTerm()
	node := &raftNode{
		id:        id,
		term:      term,
		instC:     instc,
		msgC:      msgC,
		requests:  make([]Request, 0),
		proxyReqs: make(map[uint64]Request),
		raftLog:   raftLog,
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
	r.raftLog.SaveTerm(r.term, 0)

	lastIndex, lastTerm := r.raftLog.LastIndexTerm()
	r.send(AddrPeers, &EventSolicitVoteReq{lastIndex, lastTerm})

	r.becomeRole(RoleCandidate)
}

func (r *raftNode) send(to Address, event MsgEvent) {
	msg := Message{
		from:  AddrLocal,
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
