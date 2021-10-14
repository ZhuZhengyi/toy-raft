package raft

type follower struct {
	raftNode
	leader            string
	votedFor          string
	leaderSeenTicks   uint64
	leaderSeenTimeout uint64
}

func NewFollower(node *raftNode) *follower {
	f := &follower{
		raftNode:          *node,
		leaderSeenTicks:   0,
		leaderSeenTimeout: 10,
	}

	return f
}

func (c *follower) RoleType() RoleType {
	return RoleFollower
}

func (f *follower) Step(msg *Message) {
}

func (f *follower) Tick() {
}
