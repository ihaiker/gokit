package gotcp

import (
	"time"
	"net"
)

type IdleState uint

const (
	READ_IDLE_STATE IdleState = iota //No data was received for a while.
	WRITER_IDLE_STATE                  //No data was sent for a while.
	ALL_IDLE_STATE                     //No data was either received or sent for a while.
)

//the heartbeat idle handler
type HeartbeatHandler interface {
	//handle the heartbeat event
	OnIdle(idleState IdleState, conn *Conn)
}

type HeartbeatIdle struct {
	ReadIdle  time.Duration
	WriteIdle time.Duration
	AllIdle   time.Duration
}

func DefHeartbeatIdle() *HeartbeatIdle {
	return NewHeartbeatIdle(7 * time.Second, 15 * time.Second)
}

func NewHeartbeatIdle(read, write time.Duration) *HeartbeatIdle {
	return &HeartbeatIdle{
		ReadIdle:  read,
		WriteIdle: write,
		AllIdle: (read + write),
	}
}

type Heartbeat struct {
	Handler HeartbeatHandler
	Idle    *HeartbeatIdle
}

type HeartbeatHandlerMaker (func(conn *net.TCPConn) *Heartbeat)

//Idle Package
type IdlePackage interface {
	IsIdle() bool
}