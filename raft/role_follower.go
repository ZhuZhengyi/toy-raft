package raft

type Follower struct {
	*RaftNode
	leader            string
	votedFor          string
	leaderSeenTicks   int
	leaderSeenTimeout int
}

func NewFollower(node *RaftNode) *Follower {
	f := &Follower{
		RaftNode:          node,
		leaderSeenTicks:   0,
		leaderSeenTimeout: randInt(ELECT_TICK_MIN, ELECT_TICK_MAX),
	}

	return f
}

var (
	_ RaftRole = (*Follower)(nil)
)

func (f *Follower) Type() RoleType {
	return RoleFollower
}

func (f *Follower) Step(msg *Message) {

}

func (f *Follower) Tick() {
	f.leaderSeenTicks += 1
	if f.leaderSeenTicks >= f.leaderSeenTimeout {
		f.leaderSeenTicks = 0
		f.becomeCandidate()
	}
}
