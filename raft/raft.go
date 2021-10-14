// raft.go

package raft

type raft struct {
	id          uint64
	tickc       chan struct{}
	stopc       chan struct{}
	clientc     chan Request
	peerc       chan Message
	msgc        chan Message
	node        *raftNode
	reqSessions map[ReqId]Session
}

func NewRaft(id uint64) *raft {
	r := &raft{
		id:          id,
		node:        newRaftNode(id, RoleFollower),
		reqSessions: make(map[ReqId]Session),
	}
	return r
}

func (r *raft) run() {
	for {
		select {
		case <-r.stopc:
			break
		case <-r.tickc:
			r.node.Tick()
		case msg := <-r.peerc:
			r.node.Step(&msg)
		case req := <-r.clientc:
			msg := r.getReqMsg(req)
			r.node.Step(msg)
		case msg := <-r.msgc:
			r.dispatch(msg)
		}
	}
}

func (r *raft) dispatch(msg Message) {
	switch msg.to {
	case AddrPeer, AddrPeers:
		r.peerc <- msg
	case AddrClient:
		if msg.Type() == MsgClientResp {
			r.ReplyToClient(msg.event.(*EventClientResp))
		}
	default:
	}
}

func (r *raft) pushReqSession(id ReqId, s Session) {
	r.reqSessions[id] = s
}

func (r *raft) popReqSession(id ReqId) (s Session) {
	s = r.reqSessions[id]
	delete(r.reqSessions, id)
	return
}

func (r *raft) ReplyToClient(resp *EventClientResp) {
	s := r.popReqSession(resp.id)
	s.Send(resp.response)
}

func (r *raft) getReqMsg(req Request) (msg *Message) {
	clientReq := NewEventClientReq(req)
	r.pushReqSession(clientReq.id, req.Session())
	msg = NewMessage(AddrClient, AddrLocal, 0, clientReq)
}
