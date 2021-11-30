package raft

import (
	"fmt"
	"sync/atomic"
)

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

func (f *Follower) String() string {
	return fmt.Sprintf("{id: %v, term: %v, role: %v, leader: %v, leaderSeenTicks: %v}",
		f.id, f.term, f.RoleType(), f.leader, f.leaderSeenTicks)
}

//Step a message with rsm
func (f *Follower) Step(msg *Message) {
	// 收到leader消息，重置心跳计数
	if f.isFromLeader(msg.from) {
		atomic.StoreInt64(&f.leaderSeenTicks, 0)
	}

	switch msg.MsgType() {
	case MsgTypeHeartbeatReq:
		if f.isFromLeader(msg.from) {
			hbReq := msg.event.(*EventHeartbeatReq)
			hasCommitted := f.log.Has(hbReq.commitIndex, hbReq.commitTerm)
			if hasCommitted && hbReq.commitIndex > f.log.CommittedIndex() {
				oldCommittedIndex := f.log.CommittedIndex()
				f.log.Commit(hbReq.commitIndex)
				for i := oldCommittedIndex + 1; i < hbReq.commitIndex; i++ {
					entry := f.log.Get(i)
					instruction := &InstApply{entry}
					f.instC <- instruction
				}
			}
		}
	case MsgTypeVoteReq:
		//TODO:
	case MsgTypeAppendEntriesReq:
		//TODO:
	case MsgTypeClientReq:
		//TODO:
	case MsgTypeClientResp:
		//TODO:
	case MsgTypeVoteResp:
		//TODO:
	default:
		logger.Warn("role(%v) received unexpected message(%v)\n", f, msg)
	}
}

//Tick tick
func (f *Follower) Tick() {
	atomic.AddInt64(&f.leaderSeenTicks, 1)
	if atomic.LoadInt64(&f.leaderSeenTicks) >= f.leaderSeenTimeout {
		logger.Info("%v elect tick timeout, becomeCandidate\n", f)
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
