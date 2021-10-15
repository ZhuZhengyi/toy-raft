// raft_test.go

package raft

import "testing"

func buildRaft(t *testing.T) {
	r1 := NewRaft(1, []uint64{2, 3})
	r1.Serve()
}
