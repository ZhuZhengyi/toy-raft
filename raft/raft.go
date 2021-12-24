// raft.go

package raft

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type reqSession struct {
	req     Request
	session Session
}

var (
	stopChanSignal = struct{}{}
)

//
type raft struct {
	config            *RaftConfig     //
	id                uint64          // raft id
	ticker            *time.Ticker    //
	stopc             chan struct{}   // stop signal chan
	clientInC         chan reqSession // request recv from client
	peerInC           chan *Message   // msg chan  recv from peer
	peerOutC          chan *Message   // msg send to peer
	peerSessionsMutex sync.RWMutex
	peerSessions      map[uint64]net.Conn //
	peerListener      net.Listener        //
	node              *RaftNode           //
	smDriver          *InstDriver         //
	reqSessions       map[ReqId]Session   //
	cancelFunc        context.CancelFunc  //
}

func (r *raft) String() string {
	return fmt.Sprintf("{node:%v}", r.node)
}

//NewRaft allocate a new raft struct from heap and init it
//return raft struct pointer
func NewRaft(config *RaftConfig, logStore LogStore, sm InstStateMachine) *raft {
	clientInC := make(chan reqSession, 64)
	instC := make(chan Instruction, 64)
	msgC := make(chan *Message, 64)
	peerInC := make(chan *Message, 64)
	peerOutC := make(chan *Message, 64)
	node := NewRaftNode(uint64(config.Id), RoleFollower, logStore, instC, msgC)
	instDriver := NewInstDriver(instC, msgC, sm)
	peerSessions := make(map[uint64]net.Conn)

	r := &raft{
		config:       config,
		id:           uint64(config.Id),
		stopc:        make(chan struct{}),
		clientInC:    clientInC,
		peerInC:      peerInC,
		peerOutC:     peerOutC,
		peerSessions: peerSessions,
		node:         node,
		smDriver:     instDriver,
		ticker:       time.NewTicker(TICK_INTERVAL_MS),
		reqSessions:  make(map[ReqId]Session),
	}

	return r
}

//Serve serve raft engine
func (r *raft) Start() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	r.cancelFunc = cancelFunc

	go r.smDriver.Run(ctx)
	go r.runPeerMsgOut(ctx)
	go r.runPeerMsgIn(ctx)
	go r.run(ctx)
}

//Stop stop raft
func (r *raft) Stop() {
	//r.cancelFunc()
}

//Query query for read req raft
func (r *raft) Query(command []byte) (resp []byte, err error) {

	return
}

//Mutate
func (r *raft) Mutate(command []byte) (resp []byte, err error) {

	return
}

func (r *raft) Status(rid uint64) (resp []byte, err error) {

	return
}

func (r *raft) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("raft %v run stopped", r)
			break
		case <-r.ticker.C:
			r.node.Tick()
		case msg := <-r.peerInC:
			r.node.Step(msg)
		case clientReq := <-r.clientInC:
			msg := r.makeClientMsg(clientReq.req, clientReq.session)
			r.node.Step(msg)
		case msg := <-r.node.msgC:
			r.dispatchTo(msg)
		}
	}
}

// dispatchTo message
func (r *raft) dispatchTo(msg *Message) {
	switch msg.to.(type) {
	case *AddrPeer, *AddrPeers:
		r.peerOutC <- msg
	case *AddrClient:
		switch event := msg.event.(type) {
		case *EventClientResp:
			r.replyToClient(event)
		}
	default:
		logger.Warn("dispatch invalid to msg: %v\n", msg)
	}
}

//push session with reqid
func (r *raft) pushClientSession(id ReqId, s Session) {
	r.reqSessions[id] = s
}

//
func (r *raft) popClientSession(id ReqId) (s Session) {
	s = r.reqSessions[id]
	delete(r.reqSessions, id)
	return
}

//
func (r *raft) replyToClient(resp *EventClientResp) {
	if s, ok := r.reqSessions[resp.id]; ok {
		s.Reply(resp.response)
	}
}

func (r *raft) makeClientMsg(req Request, session Session) (msg *Message) {
	clientReq := NewEventClientReq(req)
	r.pushClientSession(clientReq.id, session)
	msg = NewMessage(&AddrClient{}, &AddrLocal{}, 0, clientReq)
	return
}

func (r *raft) UpdatePeerSessions(id uint64, session net.Conn) {
	r.peerSessionsMutex.Lock()
	r.peerSessions[id] = session
	r.peerSessionsMutex.Unlock()

	r.node.peers = append(r.node.peers, id)
}

func (r *raft) GetPeerSession(id uint64) (session net.Conn, ok bool) {
	r.peerSessionsMutex.RLock()
	session, ok = r.peerSessions[id]
	r.peerSessionsMutex.RUnlock()
	return
}
