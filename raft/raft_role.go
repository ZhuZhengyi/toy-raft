// raft_role.go
package raft

//RaftRole raft node role
type RaftRole interface {
	Type() RoleType
	Step(msg *Message)
	Tick()
}

//RoleType raft role type
type RoleType int8

const (
	//RoleFollower follower role
	RoleFollower RoleType = iota
	//RoleCandidate cadidate role
	RoleCandidate
	//RoleLeader leader role
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
