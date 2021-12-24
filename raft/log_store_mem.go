// log_store_mem.go

package raft

import (
	"errors"
	"sync"
)

var (
	ErrNoEntry = errors.New("no entry exist")
)

type SyncEntries struct {
	sync.RWMutex
	entries []Entry
}

type memLogStore struct {
	se           SyncEntries
	metaMutex    sync.RWMutex
	appliedIndex uint64
	metaData     map[LogMetaKey][]byte
}

var (
	_ (LogStore) = new(memLogStore)
)

func NewMemLogStore() *memLogStore {
	return &memLogStore{
		se:       SyncEntries{entries: make([]Entry, 4096)},
		metaData: make(map[LogMetaKey][]byte),
	}
}

func (log *memLogStore) Append(entries []Entry) {
	log.Lock()
	log.entries = append(log.entries, entries...)
	log.Unlock()
}

func (log *memLogStore) AppliedIndex() uint64 {
	return log.appliedIndex
}

func (log *memLogStore) Get(index uint64) *Entry {
	log.RLock()
	defer log.RUnlock()
	for _, e := range log.entries {
		if e.index == index {
			return &e
		}
	}
	return nil
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

func (log *memLogStore) StoreMetaData(key LogMetaKey, data []byte) {
	log.metaMutex.Lock()
	log.metaData[key] = data
	log.metaMutex.Unlock()
}

func (log *memLogStore) LoadMetaData(key LogMetaKey) []byte {
	log.metaMutex.RLock()
	data := log.metaData[key]
	log.metaMutex.RUnlock()

	return data
}
