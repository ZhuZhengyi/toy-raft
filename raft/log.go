// log.go
package raft

import "sync"

type Entry struct {
	index   uint64
	term    uint64
	command []byte
}

type Log interface {
	Append(term uint64, command []byte) Entry
	Commit(index uint64)
	Get(index uint64) Entry
}

type memLog struct {
	sync.RWMutex
	lastIndex   uint64
	lastTerm    uint64
	commitIndex uint64
	commitTerm  uint64
	entries     []Entry
}

func NewMemLog() *memLog {
	return &memLog{
		entries: make([]Entry, 0),
	}
}

func (l *memLog) Append(term uint64, command []byte) Entry {
	e := Entry{
		index:   l.lastIndex + 1,
		term:    term,
		command: command,
	}
	l.Lock()
	l.lastIndex = l.lastIndex + 1
	l.lastTerm = term
	l.entries = append(l.entries, e)
	l.Unlock()

	return e
}

func (l *memLog) Get(index uint64) Entry {
	for _, e := range l.entries {
		if e.index == index {
			return e
		}
	}
	return Entry{}
}

func (l *memLog) Commit(index uint64) {
	e := l.Get(index)
	l.Lock()
	l.commitIndex = e.index
	l.commitTerm = e.term
	l.Unlock()

	return
}
