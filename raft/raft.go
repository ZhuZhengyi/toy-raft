// raft.go

package raft

import "time"

type reqSession struct {
	req     Request
	session Session
}

type raft struct {
	id      uint64 // raft id
	peers   []uint64
	stopc   chan struct{}
	clientc chan reqSession
	peerc   chan Message
	//msgc        chan Message
	node        *raftNode
	ticker      *time.Ticker
	smDriver    *InstDriver
	reqSessions map[ReqId]Session
}

func NewRaft(id uint64, peers []uint64, logStore LogStore, sm InstStateMachine) *raft {
	instC := make(chan Instruction, 64)
	msgC := make(chan Message, 64)
	node := NewRaftNode(id, RoleFollower, logStore, instC, msgC)
	r := &raft{
		id:          id,
		stopc:       make(chan struct{}),
		clientc:     make(chan reqSession, 64),
		peerc:       make(chan Message, 64),
		node:        node,
		smDriver:    NewInstDriver(instC, msgC, sm),
		ticker:      time.NewTicker(time.Duration(TickInterval) * time.Millisecond),
		reqSessions: make(map[ReqId]Session),
	}

	return r
}

func (r *raft) Serve() {
	go r.smDriver.drive()
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
		case reqSession := <-r.clientc:
			msg := r.getReqMsg(reqSession.req, reqSession.session)
			r.node.Step(msg)
		case msg := <-r.node.msgC:
			r.dispatch(msg)
		}
	}
}

func (r *raft) dispatch(msg Message) {
	switch msg.to.Type() {
	case AddrTypePeer, AddrTypePeers:
		r.peerc <- msg
	case AddrTypeClient:
		if msg.Type() == MsgTypeClientResp {
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

func (r *raft) getReqMsg(req Request, session Session) (msg *Message) {
	clientReq := NewEventClientReq(req)
	r.pushReqSession(clientReq.id, session)
	msg = NewMessage(&AddrClient{}, &AddrLocal{}, 0, clientReq)
	return
}
