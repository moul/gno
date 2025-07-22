package auth

import (
	"fmt"
	"time"

	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// SessionAccount interface for accounts that support sessions
type SessionAccount interface {
	std.Account
	GetSession(pubKey crypto.PubKey) (SessionInfo, error)
	HasSession(pubKey crypto.PubKey) bool
}

// SessionInfo interface for session details
type SessionInfo interface {
	GetExpireAt() *time.Time
	CanTransferAmount(amount int64) bool
	ConsumeTransferCapacity(amount int64) error
	HasRealmAccess(realmPath string) bool
	CanManageSessions() bool
}

// ValidateSessionForTx validates that a session key has the necessary permissions to execute a transaction
func ValidateSessionForTx(ctx sdk.Context, acc std.Account, pubKey crypto.PubKey, tx std.Tx) sdk.Result {
	// If it's the master key, no session validation needed
	if acc.GetMasterKey() != nil && acc.GetMasterKey().GetPubKey().Equals(pubKey) {
		return sdk.Result{}
	}

	// Try to cast to SessionAccount to get session details
	sessionAcc, ok := acc.(SessionAccount)
	if !ok {
		// If not a SessionAccount, we can't validate sessions
		return sdk.Result{}
	}

	// Get the session
	session, err := sessionAcc.GetSession(pubKey)
	if err != nil {
		return abciResult(std.ErrUnauthorized(fmt.Sprintf("session not found: %v", err)))
	}

	// Check if session is expired
	if expireAt := session.GetExpireAt(); expireAt != nil && ctx.BlockTime().After(*expireAt) {
		return abciResult(std.ErrUnauthorized("session has expired"))
	}

	// Validate permissions for each message in the transaction
	for _, msg := range tx.GetMsgs() {
		if res := validateSessionForMsg(ctx, session, msg); !res.IsOK() {
			return res
		}
	}

	return sdk.Result{}
}

// validateSessionForMsg validates that a session has permission to execute a specific message
func validateSessionForMsg(ctx sdk.Context, session SessionInfo, msg std.Msg) sdk.Result {
	// Check message route to determine type
	route := msg.Route()
	
	switch route {
	case "bank":
		// For bank messages, we need to check transfer capacity
		// The actual validation will be done by checking the message type name
		msgType := msg.Type()
		
		if msgType == "send" || msgType == "multi-send" {
			// For now, we'll skip validation of specific amounts
			// This would need to be handled at a higher level where we have access to message details
			// The bank module itself should check session capacity when processing transfers
		}

	// Add other routes that need session validation here
	// For example: "vm" for contract calls, etc.
	}

	return sdk.Result{}
}

// ConsumeSessionTransferCapacity consumes transfer capacity from a session after a successful transfer
func ConsumeSessionTransferCapacity(ctx sdk.Context, acc std.Account, pubKey crypto.PubKey, amount int64) error {
	// If it's the master key, no capacity consumption needed
	if acc.GetMasterKey() != nil && acc.GetMasterKey().GetPubKey().Equals(pubKey) {
		return nil
	}

	// Try to cast to SessionAccount
	sessionAcc, ok := acc.(SessionAccount)
	if !ok {
		return nil // Not a SessionAccount, no session capacity to consume
	}

	// Get the session
	session, err := sessionAcc.GetSession(pubKey)
	if err != nil {
		return fmt.Errorf("session not found: %v", err)
	}

	// Consume the capacity
	return session.ConsumeTransferCapacity(amount)
}