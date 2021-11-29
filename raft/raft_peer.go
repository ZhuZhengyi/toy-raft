// Pack age main provides ...
package raft

import (
	"context"
	"net"
	"time"
)

func (r *raft) closePeerMsgIn() {
	r.peerListener.Close()
}

// deal with peer in connect
func (r *raft) runPeerMsgIn(ctx context.Context) {
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		conn, err2 := r.peerListener.Accept()
		if err2 != nil {
			logger.Warn("listen tcp %v, error: %v\n", r.peerListener.Addr(), err2)
			return
		}

		logger.Info("raft(%v) get connect from %v \n", r, conn.RemoteAddr())

		select {
		case <-ctx.Done():
			logger.Info("doPeerRecv stopped\n")
			return
		default:
		}

		go r.doRecvFromConn(subCtx, conn)
	}
}

func (r *raft) runPeerMsgOut(ctx context.Context) {
	connTimeout := time.Duration(10)
	//
	for _, peer := range r.config.Peers {
		for i := 0; i < PEER_CONNECT_TRY_TIMES; i++ {
			peerConn, err := net.DialTimeout("tcp", peer, connTimeout)
			if err != nil {
				logger.Error("raft: %v connect to peer(%v) error:%v\n", r, peer, err)
				time.Sleep(PEER_CONNECT_TRY_SLEEP_INTERVAL)
			} else {
				r.peerSessions[peer] = peerConn
				logger.Info("get peer(%v) connection\n", peer)
				break
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("doPeerSend stopped\n")
			return
		case msg := <-r.peerOutC:
			r.sendPeerMsg(&msg)
		}
	}
}

func (r *raft) doRecvFromConn(ctx context.Context, conn net.Conn) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("doPeerRecv stopped")
			return
		default:
		}
		msg := r.recvMsgFromPeer(ctx, conn)
		r.peerInC <- *msg
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
		r.sendMsgToPeer(msg, peer)
	case AddrTypePeers:
		addrTo := msg.to.(*AddrPeers)
		for _, peer := range addrTo.peers {
			r.sendMsgToPeer(msg, peer)
		}
	default:
		logger.Warn("sentToPeer invalid msg: %v\n", msg)
	}
}

func (r *raft) sendMsgToPeer(msg *Message, peer string) {
	msg.to = &AddrPeer{peer}
	msg.from = AddressLocal
	session := r.peerSessions[peer]

	msg.SendTo(session)

	//putMessage(msg)
}

// recv message from peer
func (r *raft) recvMsgFromPeer(ctx context.Context, conn net.Conn) *Message {
	msg := new(Message)
	msg.RecvFrom(conn)
	msg.to = AddressLocal
	peer := conn.RemoteAddr().String()
	msg.from = &AddrPeer{peer: peer}

	return msg
}
