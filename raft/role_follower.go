package raft

var (
	_ RaftNode = (*Follower)(nil)
)

type Follower struct {
	*raftNode
	leader            string
	votedFor          string
	leaderSeenTicks   int
	leaderSeenTimeout int
}

func NewFollower(r *raftNode) *Follower {
	f := &Follower{
		raftNode:          r,
		leaderSeenTicks:   0,
		leaderSeenTimeout: randInt(ELECT_TICK_MIN, ELECT_TICK_MAX),
	}

	return f
}

func (f *Follower) RoleType() RoleType {
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
