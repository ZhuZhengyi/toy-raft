package raft

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//go:generate stringer -type=AddrType  -linecomment

type AddrType int8

const (
	AddrTypeLocal  AddrType = iota // AddrLocal
	AddrTypeClient                 // AddrClient
	AddrTypePeer                   // AddrPeer
	AddrTypePeers                  // AddrPeers
)

type Address interface {
	Type() AddrType
	String() string
	Size() uint64
	Marshal() []byte
	Unmarshal([]byte) error
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

func (a *AddrLocal) Size() uint64 {
	return 1
}

func (a *AddrClient) Size() uint64 {
	return 1
}

func (a *AddrPeer) Size() uint64 {
	return 1
}

func (a *AddrPeers) Size() uint64 {
	return 1
}

func (a *AddrLocal) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, AddrTypeLocal)

	return buffer.Bytes()
}

func (a *AddrClient) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, AddrTypeClient)

	return buffer.Bytes()
}

func (a *AddrPeer) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, AddrTypePeer)

	return buffer.Bytes()
}

func (a *AddrPeers) Marshal() []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, AddrTypePeers)

	return buffer.Bytes()
}

func (a *AddrLocal) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	var ad AddrType
	if err := binary.Read(buffer, binary.BigEndian, &ad); err != nil {
		logger.Warn("unmarshal %v error: %v", a, err)
	}
	if ad != AddrTypeLocal {
		logger.Warn("unmarshal %v type error: %v", a, ad)
	}
}

func (a *AddrClient) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	var ad AddrType
	if err := binary.Read(buffer, binary.BigEndian, &ad); err != nil {
		logger.Warn("unmarshal %v error: %v", a, err)
	}
	if ad != AddrTypeClient {
		logger.Warn("unmarshal %v type error: %v", a, ad)
	}
}

func (a *AddrPeer) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	var ad AddrType
	if err := binary.Read(buffer, binary.BigEndian, &ad); err != nil {
		logger.Warn("unmarshal %v error: %v", a, err)
	}
	if ad != AddrTypePeer {
		logger.Warn("unmarshal %v type error: %v", a, ad)
	}
}

func (a *AddrPeers) Unmarshal(data []byte) {
	buffer := bytes.NewBuffer(data)

	var ad AddrType
	if err := binary.Read(buffer, binary.BigEndian, &ad); err != nil {
		logger.Warn("unmarshal %v error: %v", a, err)
	}
	if ad != AddrTypePeers {
		logger.Warn("unmarshal %v type error: %v", a, ad)
	}
}
