package raft

import "sync/atomic"

//Follower a follower raft role
type Follower struct {
	*RaftNode
	leader            string
	votedFor          string
	leaderSeenTicks   int64
	leaderSeenTimeout int64
}

var (
	_ RaftRole = (*Follower)(nil)
)

//NewFollower allocate a new raft follower struct
func NewFollower(node *RaftNode) *Follower {
	f := &Follower{
		RaftNode:          node,
		leaderSeenTicks:   0,
		leaderSeenTimeout: int64(randInt(ELECT_TICK_MIN, ELECT_TICK_MAX)),
	}

	return f
}

//Type raft role type
func (f *Follower) Type() RoleType {
	return RoleFollower
}

//Step a message with rsm
func (f *Follower) Step(msg *Message) {
	// 收到leader消息，重置心跳计数
	if f.isFromLeader(msg.from) {
		atomic.StoreInt64(&f.leaderSeenTicks, 0)
	}

	switch msg.EventType() {
	case EventTypeHeartbeatReq:
		if f.isFromLeader(msg.from) {
			//TODO:
		}
	case EventTypeVoteReq:
		//TODO:
	case EventTypeAppendEntriesReq:
		//TODO:
	case EventTypeClientReq:
		//TODO:
	case EventTypeClientResp:
		//TODO:
	case EventTypeVoteResp:
		//TODO:
	default:
		logger.Warn("role(%v) received unexpected message(%v)\n", f, msg)
	}
}

//Tick tick
func (f *Follower) Tick() {
	atomic.AddInt64(&f.leaderSeenTicks, 1)
	if atomic.LoadInt64(&f.leaderSeenTicks) >= f.leaderSeenTimeout {
		atomic.StoreInt64(&f.leaderSeenTicks, 0)
		f.becomeCandidate()
	}
}

func (f *Follower) isFromLeader(addr Address) bool {
	if addr.Type() == AddrTypePeer {
		from := addr.(*AddrPeer).peer
		if from == f.leader {
			return true
		}
	}

	return false
}
