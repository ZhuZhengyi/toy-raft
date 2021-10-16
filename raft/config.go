// config.go

package raft

import "time"

const (
	TICK_INTERVAL         = 10
	TICK_INTERVAL_MS      = time.Duration(TICK_INTERVAL) * time.Millisecond
	HB_TICK_MIN           = 8
	HB_TICK_MAX           = 15
	ELECT_TICK_MIN        = 8
	ELECT_TICK_MAX        = 15
	CLIENT_REQ_BATCH_SIZE = 4096
	QUEUED_REQ_BATCH_SIZE = 4096
)
