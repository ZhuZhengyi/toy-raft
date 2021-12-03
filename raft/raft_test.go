// raft_test.go

package raft

import (
	"flag"
	"testing"
	"time"
)

var (
	TestLogLevel = "Detail"
)

func init() {
	flag.StringVar(&TestLogLevel, "LogLevel", "Debug", "test log level")
	//flag.Parse()
}

func TestRaft(t *testing.T) {
	config1 := &RaftConfig{Id: 1, PeerTcpPort: 18551, Peers: []RaftPeer{{Id: 2, Addr: "127.0.0.1:18552"}, {Id: 3, Addr: "127.0.0.1:18553"}}}
	config2 := &RaftConfig{Id: 2, PeerTcpPort: 18552, Peers: []RaftPeer{{Id: 3, Addr: "127.0.0.1:18553"}, {Id: 1, Addr: "127.0.0.1:18551"}}}
	config3 := &RaftConfig{Id: 3, PeerTcpPort: 18553, Peers: []RaftPeer{{Id: 1, Addr: "127.0.0.1:18551"}, {Id: 2, Addr: "127.0.0.1:18552"}}}

	logger.SetLogLevel(TestLogLevel)

	r1 := NewRaft(config1, NewMemLogStore(), new(DummyInstStateMachine))
	r2 := NewRaft(config2, NewMemLogStore(), new(DummyInstStateMachine))
	r3 := NewRaft(config3, NewMemLogStore(), new(DummyInstStateMachine))

	go r1.Start()
	go r2.Start()
	go r3.Start()

	defer r1.Stop()
	defer r2.Stop()
	defer r3.Stop()

	for i := 0; i < 5; i++ {
		time.Sleep(10 * time.Second)

		if r1.node.RoleType() == RoleLeader {
			r1.node.becomeFollower(r1.node.term, 0)
		} else if r2.node.RoleType() == RoleLeader {
			r2.node.becomeFollower(r2.node.term, 0)
		} else if r3.node.RoleType() == RoleLeader {
			r3.node.becomeFollower(r3.node.term, 0)
		}
	}

}
