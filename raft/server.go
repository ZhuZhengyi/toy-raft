// server.go
package raft

type server struct {
	raft *raft
}

//NewServer allocate a new server struct
func NewServer(confPath string, logStore LogStore, sm InstStateMachine) *server {
	config, err := LoadConfig(confPath)
	if err != nil {
		return nil
	}
	//id uint64, listenPort int, peers []string, logStore LogStore, sm InstStateMachine
	return &server{
		raft: NewRaft(&config.Raft, logStore, sm),
	}
}

func (s *server) Start() {
	s.raft.Serve()
}

func (s *server) Stop() {
	s.raft.Stop()
}
