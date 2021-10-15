// server.go
package raft

type server struct {
	raft *raft
}

func NewServer(id uint64, peers []uint64) *server {
	return &server{
		raft: NewRaft(id, peers),
	}
}

func (s *server) Start() {
	s.raft.Serve()
}

func (s *server) Stop() {
	s.raft.Stop()
}
