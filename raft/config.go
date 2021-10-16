// config.go

package raft

import "time"

const (
	TICK_INTERVAL         = 10
	TICK_INTERVAL_MS      = time.Duration(TICK_INTERVAL) * time.Millisecond
	ELECT_TICK_MIN        = 8
	ELECT_TICK_MAX        = 15
	CLIENT_REQ_BATCH_SIZE = 4096
)
