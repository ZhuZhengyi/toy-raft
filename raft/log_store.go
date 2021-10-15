// log.go
package raft

import (
	"fmt"
	"sync"
)

type EntryKey struct {
	index uint64
	term  uint64
}

type Entry struct {
	EntryKey
	command []byte
}

func (e *Entry) String() string {
	return fmt.Sprintf("%v %v %v", e.index, e.term, e.command)
}

type LogStore interface {
	Append(term uint64, command []byte) Entry
	Commit(index uint64)
	Get(index uint64) Entry
	LastIndexTerm() (index, term uint64)
	CommittedIndexTerm() (index, term uint64)
}

type memLogStore struct {
	sync.RWMutex
	commitIndex uint64
	commitTerm  uint64
	entries     []Entry
}

func NewMemLogStore() *memLogStore {
	return &memLogStore{
		entries: make([]Entry, 0),
	}
}

func (log *memLogStore) Append(term uint64, command []byte) Entry {
	lastIndex := uint64(0)
	log.Lock()
	le := len(log.entries)
	if le > 0 {
		lastEntry := log.entries[le-1]
		lastIndex = lastEntry.index
	}
	e := Entry{
		EntryKey: EntryKey{
			index: lastIndex + 1,
			term:  term,
		},
		command: command,
	}
	log.entries = append(log.entries, e)
	log.Unlock()

	return e
}

func (log *memLogStore) Get(index uint64) Entry {
	log.RLock()
	defer log.RUnlock()
	for _, e := range log.entries {
		if e.index == index {
			return e
		}
	}
	return Entry{}
}

func (l *memLogStore) Commit(index uint64) {
	e := l.Get(index)
	l.Lock()
	l.commitIndex = e.index
	l.commitTerm = e.term
	l.Unlock()

	return
}

func (log *memLogStore) LastIndexTerm() (index, term uint64) {
	log.RLock()
	defer log.RUnlock()
	l := len(log.entries)
	if l == 0 {
		return 0, 0
	}
	lastEntry := log.entries[l-1]
	return lastEntry.index, lastEntry.term
}

func (log *memLogStore) CommittedIndexTerm() (index, term uint64) {
	log.RLock()
	i := log.commitIndex
	t := log.commitTerm
	log.RUnlock()

	return i, t
}
