// raft_log.go
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

type LogMetaKey uint32

const (
	MetaKeyTermVoteFor LogMetaKey = iota
	MetaKeyAppliedIndex
)

//  |        store         |
//  | ---------------|-----|
//                   ^
//                applied
type LogStore interface {
	Append([]Entry)                      // append entries to store
	Get(index uint64) Entry              //
	AppliedIndex() uint64                // get applied entry index from store
	LastIndexTerm() (index, term uint64) //
	StoreMetaData(key LogMetaKey, data []byte)
	LoadMetaData(key LogMetaKey) []byte
}

//  | ----- store  -----   |
//  |          | ------- entries ---- |
//  |          | unapplied | uncommit |
//             ^           ^          ^
//           applied    committed    last
type RaftLog struct {
	sync.RWMutex
	appliedIndex uint64
	uncommit     uint64 // uncommit offset in entries
	entries      []Entry
	logStore     LogStore
}

func NewRaftLog(store LogStore) *RaftLog {
	return &RaftLog{
		appliedIndex: store.AppliedIndex(),
		uncommit:     0,
		entries:      make([]Entry, 4096),
		logStore:     store,
	}
}

func (log *RaftLog) LastIndexTerm() (index uint64, term uint64) {
	log.RLock()
	l := len(log.entries)
	if l > 0 {
		last := log.entries[l-1]
		index = last.index
		term = last.term
		log.RUnlock()
		return
	}
	log.RUnlock()

	return log.logStore.LastIndexTerm()
}

func (log *RaftLog) CommittedIndexTerm() (index uint64, term uint64) {
	log.RLock()
	defer log.RUnlock()
	if log.uncommit > 0 {
		last := log.entries[log.uncommit-1]
		index = last.index
		term = last.term
		return
	}
	return log.logStore.LastIndexTerm()
}

func (log *RaftLog) LoadTerm() (term uint64, vote uint64) {
	return
}

func (log *RaftLog) SaveTerm(term uint64, vote uint64) {
}

// append command to raft log
func (log *RaftLog) Append(term uint64, command []byte) Entry {
	log.Lock()
	defer log.Unlock()
	lastIndex, _ := log.LastIndexTerm()
	entry := Entry{
		EntryKey: EntryKey{
			index: lastIndex + 1,
			term:  term,
		},
		command: command,
	}
	log.entries = append(log.entries, entry)

	return entry
}

// commit entries which < index
func (log *RaftLog) Commit(index uint64) {
	committedIndex, _ := log.CommittedIndexTerm()
	if index <= committedIndex {
		//TODO: error
		return
	}

	log.Lock()
	defer log.Unlock()
	commitEntries := make([]Entry, 0)
	uncommit := log.uncommit
	for _, entry := range log.entries[log.uncommit:] {
		if entry.index < index {
			uncommit += 1
			commitEntries = append(commitEntries, entry)
		}
	}
	log.logStore.Append(commitEntries)
	log.uncommit = uncommit
}

//
func (log *RaftLog) Apply(index uint64) {
	log.Lock()
	defer log.Unlock()
	applying := uint64(0)
	for _, entry := range log.entries[:log.uncommit] {
		if entry.index <= index {
			applying += 1
		}
	}
	log.uncommit -= applying
	log.entries = log.entries[applying+1:]

	//log.store.SaveMetaData()
}
