// Pack age main provides ...
package raft

import (
	"context"
	"fmt"
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.config.PeerTcpPort))
	if err != nil {
		logger.Fatal("listen tcp %v", r.config.PeerTcpPort)
		return
	}
	r.peerListener = listener
	logger.Info("raft:%v listen at %v \n", r, r.peerListener.Addr())

	defer r.peerListener.Close()

	for {
		conn, err2 := r.peerListener.Accept()
		if err2 != nil {
			logger.Warn("r:%v listen tcp:%v, error:%v\n", r, r.peerListener.Addr(), err2)
			return
		}

		logger.Info("raft:%v connect in:  %v <- %v \n", r, conn.LocalAddr(), conn.RemoteAddr())

		select {
		case <-ctx.Done():
			logger.Info("doPeerRecv stopped\n")
			return
		default:
		}

		r.UpdatePeerSessions(conn)

		go r.doRecvFromConn(subCtx, conn)
	}
}

func (r *raft) runPeerMsgOut(ctx context.Context) {
	for _, peer := range r.config.Peers {
		for i := 0; i < PEER_CONNECT_TRY_TIMES; i++ {
			peerConn, err := net.DialTimeout("tcp", peer, PEER_CONNECT_TIMEOUT)
			if err != nil {
				logger.Warn("raft:%v connect out:%v error:%v", r, peer, err)
				time.Sleep(PEER_CONNECT_TRY_SLEEP_INTERVAL)
			} else {
				r.UpdatePeerSessions(peerConn)
				r.node.peers = append(r.node.peers, peer)
				logger.Info("raft:%v connect out: %v -> %v", r, peerConn.LocalAddr(), peerConn.RemoteAddr())
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
			r.sendPeerMsg(msg)
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
		r.peerInC <- msg
	}

}

func (r *raft) sendPeerMsg(msg *Message) {
	if msg.to.Type() != AddrTypePeer && msg.to.Type() != AddrTypePeers {
		logger.Warn("sentToPeer invalid msg: %v\n", msg)
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

	if session, ok := r.GetPeerSession(peer); !ok {
		logger.Warn("raft:%v session not exist for peer:%v", r, peer)
		return
	} else {
		msg.SendTo(session)
		logger.Debug("raft:%v send msg:%v", r, msg)
	}
}

// recv message from peer
func (r *raft) recvMsgFromPeer(ctx context.Context, conn net.Conn) *Message {
	msg := new(Message)
	msg.RecvFrom(conn)
	msg.to = AddressLocal
	peer := conn.RemoteAddr().String()
	msg.from = &AddrPeer{peer: peer}

	logger.Debug("recv msg %v", msg)

	return msg
}
