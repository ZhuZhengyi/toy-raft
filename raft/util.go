// util.go

package raft

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
