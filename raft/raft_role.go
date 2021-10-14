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
	role      RaftRole
	log       Log
	requests  []Request
	proxyReqs map[uint64]Request
}

var (
	_ RaftRole = (*raftNode)(nil)
)

func newRaftNode(id uint64, role RoleType) *raftNode {
	node := &raftNode{
		id:        id,
		requests:  make([]Request, 0),
		proxyReqs: make(map[uint64]Request),
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

func (node *raftNode) becomeRole(roleType RoleType) {
	switch roleType {
	case RoleCandidate:
		node.role = NewCandidate(node)
	case RoleFollower:
		node.role = NewFollower(node)
	case RoleLeader:
		node.role = NewLeader(node)
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
