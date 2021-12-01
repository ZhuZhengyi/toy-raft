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
func NewFollower(node *RaftNode, leader, votedFor string) *Follower {
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
		if msgSolictVoteReq, ok := msg.event.(*EventSolicitVoteReq); ok {
			if addr, ok := msg.from.(*AddrPeer); ok {
				if f.votedFor != "" && f.votedFor != addr.peer {
					logger.Detail("node:%v already vote for %v, refused to vote again %v", f, f.votedFor, addr.peer)
					return
				}
			}
			lastIndex, lastTerm := f.log.LastIndexTerm()
			if msgSolictVoteReq.lastTerm < lastTerm {
				logger.Detail("node:%v msg:%v lastTerm %v < node lastTerm: %v", f, msg, msgSolictVoteReq.lastTerm, lastTerm)
				return
			}
			if msgSolictVoteReq.lastTerm == lastTerm && msgSolictVoteReq.lastIndex < lastIndex {
				return
			}

			if addr, ok := msg.from.(*AddrPeer); ok {
				f.send(msg.from, &EventGrantVoteResp{})
				f.log.SaveTerm(f.term, addr.peer)
				f.votedFor = addr.peer
				logger.Detail("node:%v, vote for:%v", f, f.votedFor)
			}
		}
	case MsgTypeAppendEntriesReq:
		//TODO:
	case MsgTypeClientReq:
		//TODO:
	case MsgTypeClientResp:
		//TODO:
	case MsgTypeVoteResp:
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
	if addr.Type() == AddrTypePeer {
		from := addr.(*AddrPeer).peer
		if from == f.leader {
			return true
		}
	}

	return false
}
