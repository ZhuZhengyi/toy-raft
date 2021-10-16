package raft

var (
	_ RaftNode = (*Follower)(nil)
)

type Follower struct {
	*raftNode
	leader            string
	votedFor          string
	leaderSeenTicks   uint64
	leaderSeenTimeout uint64
}

func NewFollower(r *raftNode) *Follower {
	f := &Follower{
		raftNode:          r,
		leaderSeenTicks:   0,
		leaderSeenTimeout: 10,
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
		f.becomeCandidate()
	}
}
