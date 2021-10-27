// raft.go

package raft

import (
	"fmt"
	"net"
	"time"
)

type reqSession struct {
	req     Request
	session Session
}

var (
	stopChanSignal = struct{}{}
)

type raft struct {
	config         *RaftConfig
	id             uint64 // raft id
	listener       net.Listener
	stopc          chan struct{}       // stop signal chan
	clientInC      chan reqSession     // request recv from client
	peerInC        chan Message        // msg chan  recv from peer
	peerOutC       chan Message        // msg send to peer
	peerOutSession map[string]net.Conn //
	transport      Transport           //
	node           *RaftNode
	ticker         *time.Ticker
	smDriver       *InstDriver
	reqSessions    map[ReqId]Session
}

//NewRaft allocate a new raft struct from heap and init it
//return raft struct pointer
func NewRaft(config *RaftConfig, logStore LogStore, sm InstStateMachine) *raft {
	instC := make(chan Instruction, 64)
	msgC := make(chan Message, 64)
	peerInC := make(chan Message, 64)
	peerOutC := make(chan Message, 64)
	clientInC := make(chan reqSession, CLIENT_REQ_BATCH_SIZE)
	node := NewRaftNode(uint64(config.ID), RoleFollower, logStore, instC, msgC)
	instDriver := NewInstDriver(instC, msgC, sm)
	peerSessions := make(map[string]net.Conn)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.PeerTcpPort))
	if err != nil {
		logger.Fatal("listen tcp %v", config.PeerTcpPort)
		return nil
	}

	r := &raft{
		config:         config,
		id:             uint64(config.ID),
		stopc:          make(chan struct{}),
		listener:       listener,
		clientInC:      clientInC,
		peerInC:        peerInC,
		peerOutC:       peerOutC,
		peerOutSession: peerSessions,
		node:           node,
		smDriver:       instDriver,
		ticker:         time.NewTicker(TICK_INTERVAL_MS),
		reqSessions:    make(map[ReqId]Session),
	}

	return r
}

func (r *raft) Serve() {
	go r.smDriver.run()
	go r.runPeerMsgOut()
	go r.runPeerMsgIn(r.stopc)
	go r.run()
}

func (r *raft) Stop() {
	r.stopc <- struct{}{}

	r.smDriver.stop()
}

func (r *raft) run() {
	for {
		select {
		case <-r.stopc:
			logger.Info("raft %v run stopped", r)
			break
		case <-r.ticker.C:
			r.node.Tick()
		case msg := <-r.peerInC:
			r.node.Step(&msg)
		case reqSession := <-r.clientInC:
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
		r.peerOutC <- msg
	case AddrTypeClient:
		if msg.EventType() == EventTypeClientResp {
			r.replyToClient(msg.event.(*EventClientResp))
		}
	default:
		logger.Warn("dispatch invalid to msg: %v\n", msg)
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
