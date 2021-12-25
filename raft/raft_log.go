// raft_log.go
package raft

//LogMetaKey log meta key
type LogMetaKey uint32

const (
	MetaKeyTermVoteFor LogMetaKey = iota
	MetaKeyAppliedIndex
)

type LogScanRange struct {
	start []byte
	end   []byte
}

type LogIter interface {
	Next()
}

//LogStore log store interface
//  |        store         |
//  | ---------------|-----|
//                   ^
//                applied
type LogStore interface {
	Append(entries ...Entry) uint64          // append entries to store
	Commit(index uint64)                     // commit log entry up to index
	Committed() uint64                       //
	Get(index uint64) *Entry                 //
	Truncate(index uint64) uint64            //
	GetMetaData(key LogMetaKey) []byte       //
	SetMetaData(key LogMetaKey, data []byte) //
	IsEmpty() bool                           //
	Size() uint64                            //
	Len() uint64                             //
	Scan(start, stop []byte) LogIter         //
}

//RaftLog raft log
//  | ----- store  -----   |
//  |          | ------- entries ---- |
//  |          | unapplied | uncommit |
//             ^           ^          ^
//           applied    committed    last
type RaftLog struct {
	store       LogStore
	lastIndex   uint64
	lastTerm    uint64
	commitIndex uint64
	commitTerm  uint64
}

func NewRaftLog(store LogStore) *RaftLog {
	var (
		commitIndex uint64
		commitTerm  uint64
		lastIndex   uint64
		lastTerm    uint64
	)
	index := store.Committed()
	if index > 0 {
		e := store.Get(index)
		commitIndex = e.index
		commitTerm = e.term
	}
	if l := store.Len(); l > 0 {
		e := store.Get(l - 1)
		lastIndex = e.index
		lastTerm = e.term
	}

	return &RaftLog{
		commitIndex: commitIndex,
		commitTerm:  commitTerm,
		lastIndex:   lastIndex,
		lastTerm:    lastTerm,
		store:       store,
	}
}

//LoadTerm load term, leader from meta
func (log *RaftLog) LoadTerm() (term uint64, leader uint64) {
	return
}

//SaveTerm save term, leader meta into log store
func (log *RaftLog) SaveTerm(term uint64, leader uint64) {
	//TODO:
}

func (log *RaftLog) Get(index uint64) *Entry {
	return log.store.Get(index)
}

//Has has entry with index, term in raftlog
func (log *RaftLog) Has(index, term uint64) bool {
	entry := log.Get(index)
	if entry != nil && entry.term == term {
		return true
	}

	return false
}

// append command to raft log
func (log *RaftLog) Append(term uint64, command []byte) *Entry {
	entry := Entry{
		index:   log.lastIndex + 1,
		term:    term,
		command: command,
	}
	log.store.Append(entry)
	log.lastIndex = entry.index
	log.lastTerm = entry.term

	return &entry
}

// commit entries which < index
func (log *RaftLog) Commit(index uint64) {
	entry := log.Get(index)
	if entry == nil {
		logger.Warn("log:%v commit index entry is nil", log)
		return
	}
	log.store.Commit(index)
	log.commitIndex = entry.index
	log.commitTerm = entry.term
}

//

//Splice:
func (log *RaftLog) Splice(entries []Entry) (uint64, error) {

	return 0, nil
}

func (log *RaftLog) Truncate(index uint64) uint64 {
	var (
		truncateIndex uint64
		truncateTerm  uint64
	)
	i := log.store.Truncate(index)
	if i > 0 {
		e := log.store.Get(i)
		truncateIndex = e.index
		truncateTerm = e.term
	}

	log.lastIndex = truncateIndex
	log.lastTerm = truncateTerm

	return index
}

func (log *RaftLog) LastIndex() uint64 {
	return log.lastIndex
}

func (log *RaftLog) LastTerm() uint64 {
	return log.lastTerm
}

func (log *RaftLog) LastIndexTerm() (uint64, uint64) {
	return log.lastIndex, log.lastTerm
}

func (log *RaftLog) CommittedIndex() uint64 {
	return log.commitIndex
}

func (log *RaftLog) CommittedTerm() uint64 {
	return log.commitTerm
}
