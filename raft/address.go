package raft

type AddrType int32

type Address interface {
	Type() AddrType
	String() string
}

const (
	AddrTypeUnkown AddrType = -1
	AddrTypeLocal  AddrType = iota
	AddrTypeClient
	AddrTypePeer
	AddrTypePeers

	AddrNameUnkown = "AddrUnkown"
	AddrNameLocal  = "AddrLocal"
	AddrNameClient = "AddrClient"
	AddrNamePeer   = "AddrPeer"
	AddrNamePeers  = "AddrPeers"
)

type address struct {
}

func (a *address) Type() AddrType {
	switch Address(a).(type) {
	case *AddrLocal:
		return AddrTypeLocal
	case *AddrClient:
		return AddrTypeClient
	case *AddrPeer:
		return AddrTypePeer
	case *AddrPeers:
		return AddrTypePeers
	default:
		return AddrTypeUnkown
	}
	return AddrTypeUnkown
}

func (a *address) String() string {
	switch Address(a).(type) {
	case *AddrLocal:
		return AddrNameLocal
	case *AddrClient:
		return AddrNameClient
	case *AddrPeer:
		return AddrNamePeer
	case *AddrPeers:
		return AddrNamePeers
	default:
		return AddrNameUnkown
	}
}

type AddrLocal struct {
	address
}

type AddrClient struct {
	address
}

type AddrPeer struct {
	address
	peer uint64
}

type AddrPeers struct {
	address
	peers []uint64
}
