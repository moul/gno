package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// MockMsgSend represents a bank send message for testing
type MockMsgSend struct {
	FromAddress crypto.Address
	ToAddress   crypto.Address
	Amount      std.Coins
}

func (m MockMsgSend) Route() string                      { return "bank" }
func (m MockMsgSend) Type() string                       { return "send" }
func (m MockMsgSend) ValidateBasic() error               { return nil }
func (m MockMsgSend) GetSignBytes() []byte               { return nil }
func (m MockMsgSend) GetSigners() []crypto.Address       { return []crypto.Address{m.FromAddress} }

// TestVMKeeperSudoBankMessage tests that non-VM messages can be queued via sudo
func TestVMKeeperSudoBankMessage(t *testing.T) {
	env := setupTestEnv()
	ctx := env.vmk.MakeGnoTransactionStore(env.ctx)

	// Give "addr1" some gnots.
	addr := crypto.AddressFromPreimage([]byte("addr1"))
	acc := env.acck.NewAccountWithAddress(ctx, addr)
	env.acck.SetAccount(ctx, acc)
	env.bankk.SetCoins(ctx, addr, std.MustParseCoins(coinsString))

	// Create a realm that will queue a bank message via sudo
	const realmPath = "gno.land/r/test/sudo_bank"
	realmFiles := []*std.MemFile{
		{Name: "bank.gno", Body: `
package sudo_bank

// This simulates a realm that wants to send tokens via sudo
func SendTokens() string {
	// In a real implementation, this would call realm.Sudo(msg)
	// where msg is a MsgSend to transfer tokens
	return "would send tokens via sudo"
}
`},
		{Name: "gnomod.toml", Body: gnolang.GenGnoModLatest(realmPath)},
	}

	msg1 := NewMsgAddPackage(addr, realmPath, realmFiles)
	err := env.vmk.AddPackage(ctx, msg1)
	assert.NoError(t, err)

	// Simulate the realm queuing a bank message
	realmAddr := crypto.AddressFromPreimage([]byte("realm"))
	recipientAddr := crypto.AddressFromPreimage([]byte("recipient"))
	
	bankMsg := MockMsgSend{
		FromAddress: realmAddr,
		ToAddress:   recipientAddr,
		Amount:      std.MustParseCoins("100ugnot"),
	}

	// Queue the bank message as if realm.Sudo was called
	env.vmk.QueueSudoMessage(realmAddr, bankMsg)

	// Verify the message was queued
	sudoMessages := env.vmk.GetSudoMessages()
	require.Len(t, sudoMessages, 1)
	assert.Equal(t, "bank", sudoMessages[0].Message.(MockMsgSend).Route())
	assert.Equal(t, recipientAddr, sudoMessages[0].Message.(MockMsgSend).ToAddress)
}

// TestVMKeeperSudoMultipleMessageTypes tests queuing different message types
func TestVMKeeperSudoMultipleMessageTypes(t *testing.T) {
	env := setupTestEnv()
	
	addr := crypto.AddressFromPreimage([]byte("addr1"))
	realmAddr := crypto.AddressFromPreimage([]byte("realm1"))
	
	// Queue a VM message
	vmMsg := NewMsgCall(addr, nil, "gno.land/r/test", "Func", []string{})
	env.vmk.QueueSudoMessage(realmAddr, vmMsg)
	
	// Queue a bank message
	bankMsg := MockMsgSend{
		FromAddress: addr,
		ToAddress:   crypto.AddressFromPreimage([]byte("recipient")),
		Amount:      std.MustParseCoins("50ugnot"),
	}
	env.vmk.QueueSudoMessage(realmAddr, bankMsg)
	
	// Queue another VM message
	vmMsg2 := NewMsgAddPackage(addr, "gno.land/r/new", nil)
	env.vmk.QueueSudoMessage(realmAddr, vmMsg2)
	
	// Verify all messages were queued in order
	sudoMessages := env.vmk.GetSudoMessages()
	assert.Len(t, sudoMessages, 3)
	
	// Check message types and order
	assert.Equal(t, "vm", sudoMessages[0].Message.(MsgCall).Route())
	assert.Equal(t, "bank", sudoMessages[1].Message.(MockMsgSend).Route())
	assert.Equal(t, "vm", sudoMessages[2].Message.(MsgAddPackage).Route())
	
	// Clear and verify
	env.vmk.ClearSudoMessages()
	assert.Len(t, env.vmk.GetSudoMessages(), 0)
}