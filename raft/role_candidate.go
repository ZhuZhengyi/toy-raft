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
	switch msg.event.Type() {
	case MsgTypeHeartbeatReq:
		if msg.from.Type() == AddrTypePeer {
			from := msg.from.(*AddrPeer)
			c.becomeFollower(msg.term, from.peer).Step(msg)
			return
		}
	case MsgTypeVoteResp:
		c.votedCount++
		if c.votedCount >= c.quorum() {
			node := c.becomeLeader()
			for _, queuedReq := range node.queuedReqs {
				node.Step(&Message{
					from:  queuedReq.from,
					to:    &AddrLocal{},
					term:  0,
					event: queuedReq.event,
				})
			}
		}
	case MsgTypeClientReq:
		c.queuedReqs = append(c.queuedReqs, queuedEvent{msg.from, msg.event})
	case MsgTypeClientResp:
		event := msg.event.(*EventClientResp)
		if event.response.Type() == RespTypeStatus {
		}
		delete(c.proxyReqs, event.id)
		c.send(AddressClient, &EventClientResp{event.id, event.response})
	case MsgTypeVoteReq:
		logger.Detail("node:%v receive ignore msg:%v", c, msg)
	default:
		logger.Warn("node:%v receive error msg:%v", c, msg)
	}
}

//Tick candidate tick
func (c *Candidate) Tick() {
	c.electionTicks++
	if c.electionTicks >= c.electionTimeout {
		logger.Info("node:%v elect tick timeout, becomeCandidate\n", c)
		c.becomeCandidate()
	}
}

func (node *Candidate) becomeFollower(term uint64, leader uint64) *RaftNode {

	node.abortProxyReqs()
	node.forwardToLeaderQueued(&AddrPeer{leader})

	return node.RaftNode
}
