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
	logStore  LogStore
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
	_, lastTerm := logStore.LastIndexTerm()
	node := &raftNode{
		id:        id,
		term:      lastTerm,
		instC:     instc,
		msgC:      msgC,
		requests:  make([]Request, 0),
		proxyReqs: make(map[uint64]Request),
		logStore:  logStore,
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
