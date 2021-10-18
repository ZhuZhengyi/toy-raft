// Code generated by "stringer -type=InstType"; DO NOT EDIT.

package raft

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InstTypeUnkown - -1]
	_ = x[InstTypeAbort-1]
	_ = x[InstTypeApply-2]
	_ = x[InstTypeNotify-3]
	_ = x[InstTypeQuery-4]
	_ = x[InstTypeStatus-5]
	_ = x[InstTypeVote-6]
}

const (
	_InstType_name_0 = "InstTypeUnkown"
	_InstType_name_1 = "InstTypeAbortInstTypeApplyInstTypeNotifyInstTypeQueryInstTypeStatusInstTypeVote"
)

var (
	_InstType_index_1 = [...]uint8{0, 13, 26, 40, 53, 67, 79}
)

func (i InstType) String() string {
	switch {
	case i == -1:
		return _InstType_name_0
	case 1 <= i && i <= 6:
		i -= 1
		return _InstType_name_1[_InstType_index_1[i]:_InstType_index_1[i+1]]
	default:
		return "InstType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}