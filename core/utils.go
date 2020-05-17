package core

import (
	"runtime"

	"thunderatz.org/thor/core/types"
)

// Ping does nothing, can be used to check if core is up
func (tc *ThorCore) Ping(args types.PingArgs, reply *types.PingReply) error {
	logger.Debug().Msg("Received ping request")
	reply.Success = true
	return nil
}

// Info returns information about the current running thor core and services
func (tc *ThorCore) Info(args types.InfoArgs, reply *types.InfoReply) error {
	logger.Debug().Msg("Received info request")

	reply.NGoRoutines = runtime.NumGoroutine()
	reply.Success = true

	return nil
}
