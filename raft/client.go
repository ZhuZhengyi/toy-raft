// client.go
package raft

type Request interface {
	Session() Session
}

type Response interface {
}

type Session interface {
	Send(Response)
}
