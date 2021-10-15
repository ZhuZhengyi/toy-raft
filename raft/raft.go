// raft.go

package raft

import "time"

type raft struct {
	id          uint64 // raft id
	peers       []uint64
	stopc       chan struct{}
	clientc     chan Request
	peerc       chan Message
	msgc        chan Message
	node        *raftNode
	ticker      *time.Ticker
	reqSessions map[ReqId]Session
}

func NewRaft(id uint64, peers []uint64) *raft {
	r := &raft{
		id:          id,
		stopc:       make(chan struct{}),
		clientc:     make(chan Request, 64),
		msgc:        make(chan Message, 64),
		peerc:       make(chan Message, 64),
		node:        newRaftNode(id, RoleFollower),
		ticker:      time.NewTicker(time.Duration(TickInterval) * time.Millisecond),
		reqSessions: make(map[ReqId]Session),
	}

	return r
}

func (r *raft) Serve() {
	go r.run()
}

func (r *raft) Stop() {
	r.stopc <- struct{}{}
}

func (r *raft) run() {
	for {
		select {
		case <-r.stopc:
			break
		case <-r.ticker.C:
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
			r.replyToClient(msg.event.(*EventClientResp))
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

func (r *raft) replyToClient(resp *EventClientResp) {
	s := r.popReqSession(resp.id)
	s.Send(resp.response)
}

func (r *raft) getReqMsg(req Request) (msg *Message) {
	clientReq := NewEventClientReq(req)
	r.pushReqSession(clientReq.id, req.Session())
	msg = NewMessage(AddrClient, AddrLocal, 0, clientReq)
}
