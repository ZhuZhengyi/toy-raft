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

	return &server{
		raft: NewRaft(&config.Raft, logStore, sm),
	}
}

func (s *server) Start() {
	s.raft.Start()
}

func (s *server) Stop() {
	s.raft.Stop()
}

//Query a raft command
func (s *server) Query(rid uint64, command []byte) (resp []byte, err error) {

	return
}

//Query a raft command
func (s *server) Mutate(rid uint64, command []byte) (resp []byte, err error) {

	return
}

//Query a raft command
func (s *server) Status(rid uint64) (resp []byte, err error) {

	return
}
