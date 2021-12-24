// role_candidate.go

package raft

import (
	"fmt"
)

//Candidate role
type Candidate struct {
	*RaftNode
	electionTicks   int
	electionTimeout int
	votedCount      uint64
}

var (
	_ RaftRole = (*Candidate)(nil)
)

//NewCandidate allocate a new candidate role
func NewCandidate(node *RaftNode) *Candidate {
	c := &Candidate{
		RaftNode:        node,
		votedCount:      1,
		electionTicks:   0,
		electionTimeout: randInt(ELECT_TICK_MIN, ELECT_TICK_MAX),
	}

	return c
}

func (c *Candidate) String() string {
	return fmt.Sprintf("{id:%v,term:%v,role:%v,votedCount:%v,electionTicks:%v}",
		c.id, c.term, c.RoleType(), c.votedCount, c.electionTicks)
}

//Type candidate role type
func (c *Candidate) Type() RoleType {
	return RoleCandidate
}

//Step step candidate state by msg
func (c *Candidate) Step(msg *Message) {
	switch event := msg.event.(type) {
	case *EventHeartbeatReq:
		switch from := msg.from.(type) {
		case *AddrPeer:
			c.RaftNode.becomeFollower(msg.term, from.peer).Step(msg)
		default:
		}
	case *EventGrantVoteResp:
		c.votedCount++
		if c.votedCount >= c.quorum() {
			logger.Debug("candidate:%v get quorum vote:%v(%v)", c, c.votedCount, c.quorum())
			c.becomeLeader()
		}
	case *EventClientReq:
		c.queuedReqs = append(c.queuedReqs, queuedEvent{msg.from, msg.event})
	case *EventClientResp:
		switch event.response.(type) {
		case *RespStatus:
			//TODO:
		}
		delete(c.proxyReqs, event.id)
		c.send(AddressClient, &EventClientResp{event.id, event.response})
	case *EventSolicitVoteReq:
		logger.Detail("node:%v receive ignore msg:%v", c, msg)
	default:
		logger.Warn("node:%v receive error msg:%v", c, msg)
	}
}

//Tick candidate tick
func (c *Candidate) Tick() {
	c.electionTicks++
	if c.electionTicks >= c.electionTimeout {
		logger.Info("candidate:%v elect tick timeout, becomeCandidate\n", c)
		c.becomeCandidate()
	}
}

func (node *Candidate) becomeFollower(term uint64, leader uint64) *RaftNode {

	node.abortProxyReqs()
	node.forwardToLeaderQueued(&AddrPeer{leader})

	return node.RaftNode
}
