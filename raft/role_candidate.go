// role_candidate.go

package raft

type candidate struct {
	raftNode
	electionTicks   uint64
	electionTimeout uint64
	voteCount       uint64
}

var (
	_ RaftRole = (*candidate)(nil)
)

func NewCandidate(node *raftNode) *candidate {
	f := &candidate{
		raftNode:        *node,
		electionTicks:   0,
		electionTimeout: 10,
		voteCount:       1,
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
