// server.go
package raft

type server struct {
	raft *raft
}

//NewServer allocate a new server struct
func NewServer(id uint64, listenPort int, peers []string, logStore LogStore, sm InstStateMachine) *server {
	return &server{
		raft: NewRaft(id, listenPort, peers, logStore, sm),
	}
}

func (s *server) Start() {
	s.raft.Serve()
}

func (s *server) Stop() {
	s.raft.Stop()
}
