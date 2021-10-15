// raft_test.go

package raft

import "testing"

func buildRaft(t *testing.T) {
	logStore := NewMemLogStore()
	peers := []uint64{2, 3}
	r1 := NewRaft(1, peers, logStore, new(DummyInstStateMachine))
	r1.Serve()
}
