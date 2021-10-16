package raft

type AddrType int32

const (
	AddrTypeUnkown AddrType = -1
	AddrTypeLocal  AddrType = iota
	AddrTypeClient
	AddrTypePeer
	AddrTypePeers
)

func (a AddrType) String() string {
	switch a {
	case AddrTypeLocal:
		return "AddrLocal"
	case AddrTypeClient:
		return "AddrClient"
	case AddrTypePeer:
		return "AddrPeer"
	case AddrTypePeers:
		return "AddrPeers"
	default:
		return "AddrUnkown"
	}
}

type Address interface {
	Type() AddrType
}

type AddrLocal struct {
}

type AddrClient struct {
}

type AddrPeer struct {
	peer uint64
}

type AddrPeers struct {
	peers []uint64
}

func (a *AddrLocal) Type() AddrType {
	return AddrTypeLocal
}

func (a *AddrClient) Type() AddrType {
	return AddrTypeClient
}

func (a *AddrPeer) Type() AddrType {
	return AddrTypePeer
}

func (a *AddrPeers) Type() AddrType {
	return AddrTypePeers
}
