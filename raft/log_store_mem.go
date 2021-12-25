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
	appliedIndex uint64
	se           SyncEntries
	metaMutex    sync.RWMutex
	metaData     map[LogMetaKey][]byte
}

var (
	_ LogStore = (*memLogStore)(nil)
)

func NewMemLogStore() *memLogStore {
	return &memLogStore{
		se:       SyncEntries{entries: make([]Entry, 4096)},
		metaData: make(map[LogMetaKey][]byte),
	}
}

func (log *memLogStore) Append(entries ...Entry) uint64 {
	log.se.Lock()
	log.se.entries = append(log.se.entries, entries...)
	log.se.Unlock()
	return 0
}

func (log *memLogStore) AppliedIndex() uint64 {
	return log.appliedIndex
}

func (log *memLogStore) Get(index uint64) *Entry {
	log.se.RLock()
	defer log.se.RUnlock()
	for _, e := range log.se.entries {
		if e.index == index {
			return &e
		}
	}
	return nil
}

func (log *memLogStore) LastIndexTerm() (index, term uint64) {
	log.se.RLock()
	defer log.se.RUnlock()
	l := len(log.se.entries)
	if l == 0 {
		return 0, 0
	}
	lastEntry := log.se.entries[l-1]
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

func (log *memLogStore) Commit(index uint64) {
}

func (log *memLogStore) Committed() uint64 {
	return 0
}

func (log *memLogStore) Truncate(index uint64) uint64 {
	return 0
}

func (log *memLogStore) IsEmpty() bool {
	return true
}

func (log *memLogStore) Size() uint64 {
	return 0
}

func (log *memLogStore) Len() uint64 {
	return 0
}

func (log *memLogStore) Scan(start, stop []byte) LogIter {

	return nil
}

func (log *memLogStore) GetMetaData(key LogMetaKey) []byte {

	return nil
}

func (log *memLogStore) SetMetaData(key LogMetaKey, data []byte) {
	//
}
