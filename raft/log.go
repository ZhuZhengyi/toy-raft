// log.go
package raft

type Entry struct {
}

type Log interface {
	Append(e Entry)
}
