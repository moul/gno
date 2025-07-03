// Package realm provides realm-specific functionality for Gno smart contracts.
// This package can only be imported by realms (not pure packages).
package realm

// Message represents a transaction message that can be queued for execution
// after the current transaction completes.
type Message interface {
	// Route returns the route key for this message type
	Route() string
	
	// Type returns the type of the message
	Type() string
	
	// ValidateBasic performs basic validation of the message
	ValidateBasic() error
	
	// GetSignBytes returns the bytes to sign for this message
	GetSignBytes() []byte
}

// Sudo queues a message to be executed as a new transaction after the current
// transaction completes. The message will be executed with the calling realm
// as the author.
//
// This function can only be called from within a realm. The realm's address
// will automatically be set as the message author by the VMKeeper.
//
// Example usage:
//   m := msg.NewMsgCall(realmAddr, "gno.land/r/demo/board", "CreatePost", args...)
//   realm.Sudo(m)
//
// The message will be executed in a new transaction context with separate
// gas accounting.
func Sudo(msg Message) {
	// This is a native function that will be implemented by the VM
	sudo(msg)
}

// sudo is the native implementation injected by the VM
func sudo(msg Message)