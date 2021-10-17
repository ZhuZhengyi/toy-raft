package raft

import "fmt"

//go:generate stringer -type=AddrType  -linecomment

type AddrType int32

const (
	AddrTypeUnkown AddrType = -1   // AddrUnkown
	AddrTypeLocal  AddrType = iota // AddrLocal
	AddrTypeClient                 // AddrClient
	AddrTypePeer                   // AddrPeer
	AddrTypePeers                  // AddrPeers
)

type Address interface {
	Type() AddrType
	String() string
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

var (
	AddressLocal  = new(AddrLocal)
	AddressClient = new(AddrClient)
	AddressPeers  = new(AddrPeers)
)

func (a *AddrLocal) Type() AddrType {
	return AddrTypeLocal
}

func (a *AddrLocal) String() string {
	return a.Type().String()
}

func (a *AddrClient) Type() AddrType {
	return AddrTypeClient
}

func (a *AddrClient) String() string {
	return a.Type().String()
}

func (a *AddrPeer) Type() AddrType {
	return AddrTypePeer
}

func (a *AddrPeer) String() string {
	return fmt.Sprintf("%v(%v)", a.Type().String(), a.peer)
}

func (a *AddrPeers) Type() AddrType {
	return AddrTypePeers
}

func (a *AddrPeers) String() string {
	return fmt.Sprintf("%v(%v)", a.Type().String(), a.peers)
}
