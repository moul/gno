package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/sdk"
)

// TestSudoMessageWithContext tests that the context can properly store and retrieve sender addresses
func TestSudoMessageWithContext(t *testing.T) {
	env := setupTestEnv()
	ctx := env.vmk.MakeGnoTransactionStore(env.ctx)

	// Test data
	senderAddr := crypto.AddressFromPreimage([]byte("sender1"))
	
	// Create a context with a sudo sender
	sudoCtx := ctx.WithValue(sdk.ContextKeySudoSender, senderAddr)

	// Verify the sender can be retrieved from context
	retrievedSender := sudoCtx.Value(sdk.ContextKeySudoSender)
	assert.Equal(t, senderAddr, retrievedSender, "context should store and retrieve sender address")

	// Verify the original context doesn't have the sender
	originalSender := ctx.Value(sdk.ContextKeySudoSender)
	assert.Nil(t, originalSender, "original context should not have sudo sender")
}

// TestQueueSudoMessageWithCorrectTypes tests that messages are queued with proper types
func TestQueueSudoMessageWithCorrectTypes(t *testing.T) {
	env := setupTestEnv()
	
	// Test data
	senderAddr := crypto.AddressFromPreimage([]byte("sender1"))
	msg := NewMsgCall(senderAddr, nil, "gno.land/r/test", "Func", []string{})

	// Queue the message
	env.vmk.QueueSudoMessage(senderAddr, msg)

	// Retrieve and verify
	sudoMessages := env.vmk.GetSudoMessages()
	assert.Len(t, sudoMessages, 1)
	assert.Equal(t, senderAddr, sudoMessages[0].Sender)
	assert.Equal(t, msg, sudoMessages[0].Message)
	
	// Verify message type
	assert.IsType(t, MsgCall{}, sudoMessages[0].Message)
	msgCall := sudoMessages[0].Message.(MsgCall)
	assert.Equal(t, "gno.land/r/test", msgCall.PkgPath)
	assert.Equal(t, "Func", msgCall.Func)
}