// Pack age main provides ...
package raft

import (
	"context"
	"net"
	"time"
)

// deal with peer in connect
func (r *raft) runPeerMsgIn(stopC chan struct{}) {
	for {
		select {
		case <-stopC:
			logger.Info("doPeerRecv stopped")
			return
		default:
		}
		conn, err2 := r.listener.Accept()
		if err2 != nil {
			logger.Fatal("listen tcp %v, error: %v", r.listener, err2)
			return
		}
		ctx, _ := context.WithCancel(context.Background())

		go r.doRecvFromConn(ctx, conn)
	}
}

func (r *raft) runPeerMsgOut() {
	connTimeout := time.Duration(10)
	//
	for _, peer := range r.config.Peers {
		peerConn, err := net.DialTimeout("tcp", peer, connTimeout)
		if err != nil {
			logger.Error("connect peer(%v) error:%v\n", peer, err)
		} else {
			r.peerOutSession[peer] = peerConn
			logger.Info("get peer(%v) connection\n", peer)
		}
	}

	for {
		select {
		case <-r.stopc:
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
		msg := r.recvMsgFromPeer(conn)
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
	session := r.peerOutSession[peer]

	msg.SendTo(session)

	//putMessage(msg)
}

// recv message from peer
func (r *raft) recvMsgFromPeer(conn net.Conn) *Message {
	msg := new(Message)
	msg.RecvFrom(conn)
	msg.to = AddressLocal
	peer := conn.RemoteAddr().String()
	msg.from = &AddrPeer{peer: peer}

	return msg
}
