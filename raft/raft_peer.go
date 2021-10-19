// Pack age main provides ...
package raft

import (
	"fmt"
	"net"
	"time"
)

func (r *raft) doPeerRecv(listerPort int) {
	listenAddr := fmt.Sprintf(":%s", listerPort)
	listener, err1 := net.Listen("tcp", listenAddr)
	if err1 != nil {
		logger.Fatal("listen tcp %v", listenAddr)
	}

	// run peer recv in task
	go func() {
		for {
			inConn, err2 := listener.Accept()
			if err2 != nil {
				logger.Fatal("listen tcp %v, error: %v", listenAddr, err2)
				return
			}
			go func() {
				for {
					msg := r.recvFromPeer(inConn)
					r.peerInC <- *msg
				}
			}()
		}
	}()
}

func (r *raft) doPeerSend() {
	connTimeout := time.Duration(10)

	// init peer out session
	for _, peer := range r.peers {
		peerConn, err := net.DialTimeout("tcp", peer, connTimeout)
		if err != nil {
			logger.Error("")
		}
		r.peerOutSession[peer] = peerConn
	}

	// run peer send out task
	go func() {
		for {
			select {
			case msg := <-r.peerOutC:
				r.sendPeerMsg(&msg)
			}
		}
	}()
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
	session := r.peerOutSession[peer]

	msg.SendTo(session)
}

func (r *raft) recvFromPeer(conn net.Conn) *Message {
	msg := new(Message)

	msg.RecvFrom(conn)

	if msg.to.Type() != AddrTypePeer {
		logger.Warn("recv invalid peer msg: %v, type not AddrTypePeer\n", msg)
		return nil
	}

	msg.to = AddressLocal
	msg.from = &AddrPeer{}

	return msg
}
