// Pack age main provides ...
package raft

import (
	"context"
	"fmt"
	"net"
	"time"
)

type RaftPeer struct {
	Id   uint64 `yaml:"Id"`
	Addr string `yaml:"Addr"`
}

func (r *raft) closePeerMsgIn() {
	r.peerListener.Close()
}

// deal with peer in connect
func (r *raft) runPeerMsgIn(ctx context.Context) {
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.config.PeerTcpPort))
	if err != nil {
		logger.Fatal("raft:%v listen tcp %v", r, r.config.PeerTcpPort)
		return
	}
	r.peerListener = listener
	logger.Info("raft:%v listen at %v \n", r, r.peerListener.Addr())

	defer r.peerListener.Close()

	for {
		conn, err2 := r.peerListener.Accept()
		if err2 != nil {
			logger.Warn("raft:%v listen tcp:%v, error:%v\n", r, r.peerListener.Addr(), err2)
			return
		}

		logger.Info("raft:%v connect in:  %v <- %v \n", r, conn.LocalAddr(), conn.RemoteAddr())

		select {
		case <-ctx.Done():
			logger.Info("doPeerRecv stopped\n")
			return
		default:
		}

		//r.UpdatePeerSessions(conn)

		go r.doRecvFromConn(subCtx, conn)
	}
}

func (r *raft) runPeerMsgOut(ctx context.Context) {
	for _, peer := range r.config.Peers {
		for i := 0; i < PEER_CONNECT_TRY_TIMES; i++ {
			peerConn, err := net.DialTimeout("tcp", peer.Addr, PEER_CONNECT_TIMEOUT)
			if err != nil {
				logger.Warn("raft:%v connect out:%v error:%v", r, peer, err)
				time.Sleep(PEER_CONNECT_TRY_SLEEP_INTERVAL)
			} else {
				r.UpdatePeerSessions(peer.Id, peerConn)
				//r.node.peers = append(r.node.peers, peer.Id)
				logger.Info("raft:%v connect out: %v:%v -> %v:%v", r, r.id, peerConn.LocalAddr(), peer.Id, peerConn.RemoteAddr())
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
	switch msg.to.(type) {
	case *AddrPeer:
		r.sendMsgToPeer(msg)
	case *AddrPeers:
		for _, peer := range r.node.peers {
			msg.to = &AddrPeer{peer: peer}
			r.sendMsgToPeer(msg)
		}
	default:
		logger.Warn("sentToPeer invalid msg: %v\n", msg)
	}
}

func (r *raft) sendMsgToPeer(msg *Message) {
	msg.from = &AddrPeer{peer: r.id}
	to := msg.To()
	if to == 0 {
		logger.Error("node:%v msg:%v to is invalid", r, msg)
		return
	}
	if session, ok := r.GetPeerSession(msg.To()); !ok {
		logger.Warn("raft:%v session not exist for peer:%v", r, msg.To())
		return
	} else {
		msg.SendTo(session)
		logger.Detail("raft:%v send msg:%v", r, msg)
	}
}

// recv message from peer
func (r *raft) recvMsgFromPeer(ctx context.Context, conn net.Conn) *Message {
	msg := new(Message)
	msg.RecvFrom(conn)

	logger.Detail("raft:%v recv msg:%v", r, msg)

	return msg
}
