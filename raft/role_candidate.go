// role_candidate.go

package raft

type Candidate struct {
	*RaftNode
	electionTicks   int
	electionTimeout int
	votedCount      uint64
}

func NewCandidate(node *RaftNode) *Candidate {
	c := &Candidate{
		RaftNode:        node,
		votedCount:      1,
		electionTicks:   0,
		electionTimeout: randInt(ELECT_TICK_MIN, ELECT_TICK_MAX),
	}

	return c
}

var (
	_ RaftRole = (*Candidate)(nil)
)

func (c *Candidate) Type() RoleType {
	return RoleCandidate
}

func (c *Candidate) Step(msg *Message) {
	if !c.validateMsg(msg) {
		return
	}

	if msg.term > c.term {
		if msg.from.Type() == AddrTypePeer {
			from := msg.from.(*AddrPeer)
			c.becomeFollower(msg.term, from.peer).Step(msg)
			return
		}
	}

	switch msg.event.Type() {
	case EventTypeHeartbeatReq:
		if msg.from.Type() == AddrTypePeer {
			from := msg.from.(*AddrPeer)
			c.becomeFollower(msg.term, from.peer).Step(msg)
			return
		}
	case EventTypeVoteResp:
		c.votedCount += 1
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
	case EventTypeClientReq:
		c.queuedReqs = append(c.queuedReqs, queuedEvent{msg.from, msg.event})
	case EventTypeClientResp:
		event := msg.event.(*EventClientResp)
		if event.response.Type() == RespTypeStatus {
		}
		delete(c.proxyReqs, event.id)
		c.send(AddressClient, &EventClientResp{event.id, event.response})
	case EventTypeVoteReq:
	default:
		logger.Warn("RaftRole(%v) reciev error msg: (%v)\n", c, msg)
		//TODO: warn
	}
}

func (c *Candidate) Tick() {
	c.electionTicks += 1
	if c.electionTicks >= c.electionTimeout {
		c.becomeCandidate()
	}
}
