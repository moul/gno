package realm

import (
	"fmt"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
)

// SudoMessage represents a message queued by a realm for later execution
type SudoMessage struct {
	RealmAddr string      // Address of the realm that queued the message
	Message   interface{} // The actual message to execute
}

// GetSudoMessages returns the list of sudo messages queued in the current context
func GetSudoMessages(m *gno.Machine) []SudoMessage {
	// Get or initialize the sudo message queue from the machine context
	queueKey := "realm.sudo.queue"
	queueValue := m.Context.Get(queueKey)
	
	if queueValue == nil {
		return nil
	}
	
	queue, ok := queueValue.([]SudoMessage)
	if !ok {
		panic("invalid sudo message queue type")
	}
	
	return queue
}

// ClearSudoMessages clears the sudo message queue
func ClearSudoMessages(m *gno.Machine) {
	queueKey := "realm.sudo.queue"
	m.Context.Set(queueKey, nil)
}

// X_sudo implements the native realm.Sudo function
func X_sudo(m *gno.Machine, msg gno.TypedValue) {
	// Ensure this is called from a realm
	realmPath := m.Realm.Path
	if realmPath == "" {
		panic("realm.Sudo can only be called from within a realm")
	}
	
	// Get the realm address
	realmAddr := m.Realm.Addr
	if realmAddr == nil {
		panic("realm address not set")
	}
	
	// Extract the message interface
	msgValue := msg.GetInterface()
	if msgValue == nil {
		panic("message cannot be nil")
	}
	
	// Get the VMKeeper from the context
	ctx := m.Context.(stdlibs.ExecContext)
	if ctx.Banker == nil {
		panic("context does not have access to VMKeeper")
	}
	
	// The Banker interface needs to be extended or we need another way
	// to access the VMKeeper. For now, we'll use the context store.
	// This is a temporary solution - in production, we'd properly extend
	// the interfaces to support this functionality.
	
	// Store the sudo message in the context for later retrieval
	queueKey := "realm.sudo.queue"
	queueValue := ctx.ContextStore.Get(queueKey)
	
	var queue []SudoMessage
	if queueValue != nil {
		// Deserialize the queue
		queue = queueValue.([]SudoMessage)
	}
	
	// Append the new message
	queue = append(queue, SudoMessage{
		RealmAddr: realmAddr.String(),
		Message:   msgValue,
	})
	
	// Store back in context
	ctx.ContextStore.Set(queueKey, queue)
}

// ValidateSudoMessage validates that a message can be executed via Sudo
func ValidateSudoMessage(msg interface{}) error {
	// Check if message implements required interface methods
	type validator interface {
		ValidateBasic() error
	}
	
	if v, ok := msg.(validator); ok {
		return v.ValidateBasic()
	}
	
	return fmt.Errorf("message does not implement ValidateBasic")
}