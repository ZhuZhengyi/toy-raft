package raft

import "fmt"

//go:generate stringer -type=AddrType  -linecomment

type AddrType int32

const (
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
	peer string
}

type AddrPeers struct {
	peers []string
}

var (
	AddressLocal  = new(AddrLocal)
	AddressClient = new(AddrClient)
	AddressPeers  = new(AddrPeers)
)

func (a *AddrLocal) String() string {
	return a.Type().String()
}

func (a *AddrClient) String() string {
	return a.Type().String()
}

func (a *AddrPeer) String() string {
	return fmt.Sprintf("%v(%v)", a.Type().String(), a.peer)
}

func (a *AddrPeers) String() string {
	return fmt.Sprintf("%v(%v)", a.Type().String(), a.peers)
}
