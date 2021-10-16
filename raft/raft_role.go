// raft_role.go
package raft

type RoleType int8

const (
	RoleFollower RoleType = iota
	RoleCandidate
	RoleLeader
)

func (r RoleType) String() string {
	switch r {
	case RoleFollower:
		return "Follower"
	case RoleCandidate:
		return "Candidate"
	case RoleLeader:
		return "Leader"
	default:
		return "Unkown"
	}
}
