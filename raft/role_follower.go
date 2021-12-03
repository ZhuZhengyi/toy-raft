package raft

import (
	"fmt"
	"sync/atomic"
)

//Follower a follower raft role
type Follower struct {
	*RaftNode
	leader            uint64
	votedFor          uint64
	leaderSeenTicks   int64
	leaderSeenTimeout int64
}

var (
	_ RaftRole = (*Follower)(nil)
)

//NewFollower allocate a new raft follower struct
func NewFollower(node *RaftNode, leader, votedFor uint64) *Follower {
	f := &Follower{
		RaftNode:          node,
		leader:            leader,
		votedFor:          votedFor,
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
	return fmt.Sprintf("{id:%v,term:%v,role:%v,leader:%v,leaderSeenTicks:%v}",
		f.id, f.term, f.RoleType(), f.leader, f.leaderSeenTicks)
}

//Step a message with rsm
func (f *Follower) Step(msg *Message) {
	// 收到leader消息，重置心跳计数
	if f.isFromLeader(msg.from) {
		atomic.StoreInt64(&f.leaderSeenTicks, 0)
	}

	switch event := msg.event.(type) {
	case *EventHeartbeatReq:
		if f.isFromLeader(msg.from) {
			hasCommitted := f.log.Has(event.commitIndex, event.commitTerm)
			if hasCommitted && event.commitIndex > f.log.CommittedIndex() {
				oldCommittedIndex := f.log.CommittedIndex()
				f.log.Commit(event.commitIndex)
				for i := oldCommittedIndex + 1; i < event.commitIndex; i++ {
					entry := f.log.Get(i)
					instruction := &InstApply{entry}
					f.instC <- instruction
				}
			}
		}
	case *EventSolicitVoteReq:
		lastIndex, lastTerm := f.log.LastIndexTerm()
		if event.lastTerm < lastTerm {
			logger.Detail("Follower:%v msg:%v lastTerm %v < node lastTerm: %v", f, msg, event.lastTerm, lastTerm)
			return
		}
		if event.lastTerm == lastTerm && event.lastIndex < lastIndex {
			return
		}
		switch from := msg.from.(type) {
		case *AddrPeer:
			if f.votedFor != 0 && f.votedFor != from.peer {
				logger.Detail("Follower:%v already vote for %v, refused to vote again %v", f, f.votedFor, from.peer)
				return
			}
			f.send(msg.from, &EventGrantVoteResp{})
			f.log.SaveTerm(f.term, from.peer)
			f.votedFor = from.peer
			logger.Detail("Follower:%v, vote for:%v", f, f.votedFor)
		}
	case *EventAppendEntriesReq:
		//TODO:
	case *EventClientReq:
		//TODO:
	case *EventClientResp:
		//TODO:
	case *EventGrantVoteResp:
		//TODO:
	default:
		logger.Warn("node:%v received unexpected message(%v)\n", f, msg)
	}
}

//Tick tick
func (f *Follower) Tick() {
	atomic.AddInt64(&f.leaderSeenTicks, 1)
	if atomic.LoadInt64(&f.leaderSeenTicks) >= f.leaderSeenTimeout {
		logger.Info("node:%v elect tick timeout", f)
		atomic.StoreInt64(&f.leaderSeenTicks, 0)
		f.becomeCandidate()
	}
}

func (f *Follower) isFromLeader(addr Address) bool {
	switch from := addr.(type) {
	case *AddrPeer:
		if from.peer == f.leader {
			return true
		}
	default:
	}

	return false
}
