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

type raft struct {
	id           uint64 // raft id
	peers        []string
	stopc        chan struct{}
	clientInC    chan reqSession
	peerInC      chan Message
	peerOutC     chan Message
	peerSessions map[string]net.Conn
	node         *RaftNode
	ticker       *time.Ticker
	smDriver     *InstDriver
	reqSessions  map[ReqId]Session
}

func NewRaft(id uint64, peers []string, logStore LogStore, sm InstStateMachine) *raft {
	instC := make(chan Instruction, 64)
	msgC := make(chan Message, 64)
	peerInC := make(chan Message, 64)
	peerOutC := make(chan Message, 64)
	clientInC := make(chan reqSession, CLIENT_REQ_BATCH_SIZE)
	node := NewRaftNode(id, RoleFollower, logStore, instC, msgC)

	peerSessions := make(map[string]net.Conn)

	r := &raft{
		id:           id,
		stopc:        make(chan struct{}),
		clientInC:    clientInC,
		peerInC:      peerInC,
		peerOutC:     peerOutC,
		peerSessions: peerSessions,
		node:         node,
		smDriver:     NewInstDriver(instC, msgC, sm),
		ticker:       time.NewTicker(TICK_INTERVAL_MS),
		reqSessions:  make(map[ReqId]Session),
	}

	return r
}

type msgChan chan Message

func (r *raft) doPeerRecv(listerPort int, peerInC chan Message) {
	listenAddr := fmt.Sprintf(":%s", listerPort)
	listener, err1 := net.Listen("tcp", listenAddr)
	if err1 != nil {
		logger.Fatal("listen tcp %v", listenAddr)
	}

	for {
		recvConn, err2 := listener.Accept()
		if err2 != nil {
			logger.Fatal("listen tcp %v", listenAddr)
		}
		//buff := make([]byte, 10)
		go func() {
			for {
				msg := r.recvFromPeer(recvConn)
				r.peerInC <- *msg
			}
		}()
	}
}

func (r *raft) doPeerSend(peers []string, peerOutC chan Message) {
	connTimeout := time.Duration(10)

	for _, peer := range peers {
		peerConn, err := net.DialTimeout("tcp", peer, connTimeout)
		if err != nil {
			logger.Error("")
		}
		r.peerSessions[peer] = peerConn
	}

	for {
		select {
		case msg := <-peerOutC:
			r.sendPeerMsg(&msg)

			//TODO: send out
		}
	}
}

func (r *raft) sendPeerMsg(msg *Message) {
	if msg.to.Type() != AddrTypePeer || msg.to.Type() != AddrTypePeers {
		return
	}

	switch msg.to.Type() {
	case AddrTypePeer:
		addrTo := msg.to.(*AddrPeer)
		peer := addrTo.peer
		r.sendToPeer(msg, peer)
	case AddrTypePeers:
		addrTo := msg.to.(*AddrPeers)
		for _, peer := range addrTo.peers {
			r.sendToPeer(msg, peer)
		}
	default:
		logger.Warn("sentToPeer invalid msg: %v\n", msg)
	}
}

func (r *raft) sendToPeer(msg *Message, peer string) {
	msg.to = &AddrPeer{peer}
	msg.from = AddressLocal
	session := r.peerSessions[peer]

	msgBuff := make([]byte, 4096)
	session.Write(msgBuff)
}

func (r *raft) recvFromPeer(conn net.Conn) *Message {
	msg := new(Message)
	msgBuff := make([]byte, 4096)
	conn.Read(msgBuff)

	if msg.to.Type() != AddrTypePeer {
		logger.Warn("recv invalid peer msg: %v, type not AddrTypePeer\n", msg)
		return nil
	}

	msg.to = AddressLocal
	msg.from = &AddrPeer{}

	return msg
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
