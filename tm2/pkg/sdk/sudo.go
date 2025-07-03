package sdk

import "github.com/gnolang/gno/tm2/pkg/crypto"

// contextKeySudoSender is the context key for the sudo message sender
type contextKeySudoSender struct{}

// ContextKeySudoSender is used to set/get the sudo message sender in the context
var ContextKeySudoSender = contextKeySudoSender{}

// SudoKeeper is an interface that allows modules to queue messages
// for execution after the current transaction completes.
type SudoKeeper interface {
	// GetSudoMessages returns all queued sudo messages
	GetSudoMessages() []SudoMessage
	// ClearSudoMessages removes all queued sudo messages
	ClearSudoMessages()
}

// SudoMessage represents a message queued for later execution
type SudoMessage struct {
	Sender  crypto.Address // Address that queued the message
	Message Msg           // The actual message to execute
}