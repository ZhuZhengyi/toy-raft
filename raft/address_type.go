package raft

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
