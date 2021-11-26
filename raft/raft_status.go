package raft

//RaftStatus raft rsm status
type RaftStatus struct {
	Server        string
	Leader        string
	Term          uint64
	PeerLastIndex map[string]uint64
	CommitIndex   uint64
	ApplyIndex    uint64
	Storage       string
	StorageSize   uint64
}
