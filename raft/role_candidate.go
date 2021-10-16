// role_candidate.go

package raft

type candidate struct {
	*raftNode
	electionTicks   uint64
	electionTimeout uint64
	votedCount      uint64
}

var (
	_ RaftRole = (*candidate)(nil)
)

func NewCandidate(node *raftNode) *candidate {
	f := &candidate{
		raftNode:        node,
		electionTicks:   0,
		electionTimeout: 10,
		votedCount:      1,
	}

	return f
}

func (c *candidate) RoleType() RoleType {
	return RoleCandidate
}

func (c *candidate) Step(msg *Message) {
}

func (c *candidate) Tick() {
}
