package raft

type Address uint32

const (
	AddrLocal Address = iota
	AddrClient
	AddrPeer
	AddrPeers

	AddrNameUnkown = "Unkown"
	AddrNameLocal  = "Local"
	AddrNameClient = "Client"
	AddrNamePeer   = "Peer"
	AddrNamePeers  = "Peers"
)

func (a Address) String() string {
	switch a {
	case AddrLocal:
		return AddrNameLocal
	case AddrClient:
		return AddrNameClient
	case AddrPeer:
		return AddrNamePeer
	case AddrPeers:
		return AddrNamePeers
	default:
		return AddrNameUnkown
	}
}
