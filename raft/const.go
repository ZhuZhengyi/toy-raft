// config.go

package raft

import "time"

const (
	DEF_PEER_TCP_PORT = 211018
)

const (
	TICK_INTERVAL_MS                = 500 * time.Millisecond
	HB_LEADER_TIMEOUT               = 3
	HB_TICK_MIN                     = 8
	HB_TICK_MAX                     = 15
	ELECT_TICK_MIN                  = 8
	ELECT_TICK_MAX                  = 15
	PEER_CONNECT_TRY_TIMES          = 5               //
	PEER_CONNECT_TIMEOUT            = 3 * time.Second //
	PEER_CONNECT_TRY_SLEEP_INTERVAL = 3 * time.Second
	CLIENT_REQ_BATCH_SIZE           = 4096
	QUEUED_REQ_BATCH_SIZE           = 4096
)
