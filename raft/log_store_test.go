//
package raft

import "testing"

func TestMemLog(t *testing.T) {
	l := NewMemLogStore()

	command := []byte{}
	entries := []Entry{
		{EntryKey{1, 1}, command},
		{EntryKey{2, 1}, command},
		{EntryKey{3, 2}, command},
		{EntryKey{4, 2}, command},
		{EntryKey{5, 3}, command},
		{EntryKey{6, 4}, command},
		{EntryKey{7, 5}, command},
	}

	for _, e := range entries {
		l.Append(e.term, e.command)
	}

	for i, e := range entries {
		el := l.Get(uint64(i + 1))
		if el.term != e.term {
			t.Errorf("log terrm error: %v %v", el, e)
		}
	}

	lastIndex, lastTerm := l.LastIndexTerm()
	le := len(entries)
	lastEntry := entries[le-1]
	if lastIndex != lastEntry.index || lastTerm != lastEntry.term {
		t.Errorf("log term error: %v %v", lastIndex, lastTerm)
	}

}
