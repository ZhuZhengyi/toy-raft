//
package raft

import "testing"

func TestMemLog(t *testing.T) {
	l := NewMemLogStore()

	command := []byte{}
	entries := []Entry{
		{1, 1, command},
		{2, 1, command},
		{3, 2, command},
		{4, 2, command},
		{5, 3, command},
		{6, 4, command},
		{7, 5, command},
	}

	l.Append(entries)

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
