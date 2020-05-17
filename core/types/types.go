// Package types exposes types for the core rpc exposed functions
package types

type (
	// SendArgs is the arguments for the ThorCore.SendDiscordAlert function
	SendArgs struct {
		Msg string
	}

	// SendReply is the reply for the ThorCore.SendDiscordAlert function
	SendReply struct {
		Success bool
	}
)

type (
	// PingArgs is the arguments for the ThorCore.Ping function
	PingArgs struct{}

	// PingReply is the reply for the ThorCore.Ping function
	PingReply struct {
		Success bool
	}
)

type (
	// InfoArgs is the arguments for the ThorCore.Info function
	InfoArgs struct{}

	// InfoReply is the reply for the ThorCore.Info function
	InfoReply struct {
		NGoRoutines int
		Success     bool
	}
)
