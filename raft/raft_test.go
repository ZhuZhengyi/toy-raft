// raft_test.go

package raft

import "testing"

func TestRaft(t *testing.T) {
	config1 := &RaftConfig{ID: 1, PeerTcpPort: 18551, Peers: []string{"127.0.0.1:18552", "127.0.0.1:18553"}}
	config2 := &RaftConfig{ID: 2, PeerTcpPort: 18552, Peers: []string{"127.0.0.1:18553", "127.0.0.1:18551"}}
	config3 := &RaftConfig{ID: 3, PeerTcpPort: 18553, Peers: []string{"127.0.0.1:18551", "127.0.0.1:18552"}}

	r1 := NewRaft(config1, NewMemLogStore(), new(DummyInstStateMachine))
	r2 := NewRaft(config2, NewMemLogStore(), new(DummyInstStateMachine))
	r3 := NewRaft(config3, NewMemLogStore(), new(DummyInstStateMachine))

	r1.Start()
	t.Logf("%v start\n", r1)
	r2.Start()
	t.Logf("%v start\n", r2)
	r3.Start()
	t.Logf("%v start\n", r3)

	defer r1.Stop()
	defer r2.Stop()
	defer r3.Stop()

}
