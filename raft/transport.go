// transport.go

package raft

//Transport message between peers
type Transport interface {
	Listen()
	Accept()
	RecvFrom()
	SendTo()
	Close()
}
