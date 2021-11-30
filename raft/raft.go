// raft.go

package raft

import (
	"context"
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

//
type raft struct {
	config       *RaftConfig         //
	id           uint64              // raft id
	ticker       *time.Ticker        //
	stopc        chan struct{}       // stop signal chan
	clientInC    chan reqSession     // request recv from client
	peerInC      chan *Message       // msg chan  recv from peer
	peerOutC     chan *Message       // msg send to peer
	peerSessions map[string]net.Conn //
	peerListener net.Listener        //
	node         *RaftNode           //
	smDriver     *InstDriver         //
	reqSessions  map[ReqId]Session   //
	cancelFunc   context.CancelFunc  //
}

func (r *raft) String() string {
	return fmt.Sprintf("{node: %v, port: %v}", r.node, r.config.PeerTcpPort)
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
	peerSessions := make(map[string]net.Conn)

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
	//logger.Debug("dispatch msg: %v", msg)
	switch msg.to.Type() {
	case AddrTypePeer, AddrTypePeers:
		r.peerOutC <- msg
	case AddrTypeClient:
		if msg.MsgType() == MsgTypeClientResp {
			r.replyToClient(msg.event.(*EventClientResp))
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
	s := r.popClientSession(resp.id)
	s.Reply(resp.response)
}

func (r *raft) makeClientMsg(req Request, session Session) (msg *Message) {
	clientReq := NewEventClientReq(req)
	r.pushClientSession(clientReq.id, session)
	msg = NewMessage(&AddrClient{}, &AddrLocal{}, 0, clientReq)
	return
}
