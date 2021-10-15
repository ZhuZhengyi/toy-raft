// raft_role_test.go

package raft

import "testing"

func TestRaftRole(t *testing.T) {
	role := RoleFollower
	rn1 := newRaftNode(1, role)

	if rn1.RoleType() != role {
		t.Errorf("role type error, %v %v", rn1.RoleType(), role)
	}

	role = RoleCandidate
	rn1.becomeRole(role)
	if rn1.RoleType() != role {
		t.Errorf("role type error, %v %v", rn1.RoleType(), role)
	}

	role = RoleLeader
	rn1.becomeRole(role)
	if rn1.RoleType() != role {
		t.Errorf("role type error, %v %v", rn1.RoleType(), role)
	}

	rn1.Step(&Message{
		from: AddrClient,
	})

}
