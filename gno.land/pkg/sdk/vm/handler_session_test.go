package vm

import (
	"testing"
	"time"

	"github.com/gnolang/gno/gno.land/pkg/gnoland"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/crypto/ed25519"
	"github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/std"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionHandler_CreateSession(t *testing.T) {
	// Setup
	env := setupTestEnv(t)
	ctx := env.ctx
	keeper := env.keeper
	handler := NewSessionHandler(keeper)

	// Create test account
	priv := ed25519.GenPrivKey()
	addr := priv.PubKey().Address()
	acc := &gnoland.GnoAccount{
		BaseAccount: std.BaseAccount{
			Address:       addr,
			PubKey:        nil,
			AccountNumber: 1,
			SequenceSum:   0,
		},
		Sessions: []gnoland.GnoSession{},
	}
	keeper.setAccount(ctx, acc)

	// Create session key
	sessionPriv := ed25519.GenPrivKey()
	sessionPub := sessionPriv.PubKey()
	
	// Test successful session creation
	t.Run("CreateSessionSuccess", func(t *testing.T) {
		msg := MsgCreateSession{
			Signer:           addr,
			Name:             "test-session",
			PubKey:           sessionPub,
			TransferCapacity: 1000,
			RealmWhitelist:   []string{"gno.land/r/demo/*"},
			Flags:            0,
		}

		res := handler.Process(ctx, msg)
		assert.True(t, res.IsOK(), "expected successful result, got: %v", res)

		// Verify session was added
		updatedAcc := keeper.getAccount(ctx, addr).(*gnoland.GnoAccount)
		require.Len(t, updatedAcc.Sessions, 1)
		assert.Equal(t, "test-session", updatedAcc.Sessions[0].Name)
		assert.Equal(t, sessionPub, updatedAcc.Sessions[0].PubKey)
		assert.Equal(t, int64(1000), updatedAcc.Sessions[0].TransferCapacity)
	})

	// Test session with expiration
	t.Run("CreateSessionWithExpiration", func(t *testing.T) {
		sessionPriv2 := ed25519.GenPrivKey()
		sessionPub2 := sessionPriv2.PubKey()
		expireTime := ctx.BlockTime().Add(1 * time.Hour)
		
		msg := MsgCreateSession{
			Signer:           addr,
			Name:             "temp-session",
			PubKey:           sessionPub2,
			ExpireAt:         &expireTime,
			TransferCapacity: 500,
		}

		res := handler.Process(ctx, msg)
		assert.True(t, res.IsOK())

		// Verify session was added with expiration
		updatedAcc := keeper.getAccount(ctx, addr).(*gnoland.GnoAccount)
		require.Len(t, updatedAcc.Sessions, 2)
		assert.NotNil(t, updatedAcc.Sessions[1].ExpireAt)
		assert.Equal(t, expireTime, *updatedAcc.Sessions[1].ExpireAt)
	})

	// Test duplicate session
	t.Run("CreateDuplicateSession", func(t *testing.T) {
		msg := MsgCreateSession{
			Signer: addr,
			Name:   "test-session",
			PubKey: sessionPub, // Same pubkey as first session
		}

		res := handler.Process(ctx, msg)
		assert.False(t, res.IsOK())
		assert.Contains(t, res.Err, "session already exists")
	})

	// Test non-existent account
	t.Run("CreateSessionNonExistentAccount", func(t *testing.T) {
		nonExistentAddr := crypto.Address("nonexistent")
		msg := MsgCreateSession{
			Signer: nonExistentAddr,
			Name:   "session",
			PubKey: sessionPub,
		}

		res := handler.Process(ctx, msg)
		assert.False(t, res.IsOK())
		assert.Contains(t, res.Err, "does not exist")
	})
}

func TestSessionHandler_RevokeSession(t *testing.T) {
	// Setup
	env := setupTestEnv(t)
	ctx := env.ctx
	keeper := env.keeper
	handler := NewSessionHandler(keeper)

	// Create test account with sessions
	priv := ed25519.GenPrivKey()
	addr := priv.PubKey().Address()
	
	sessionPriv1 := ed25519.GenPrivKey()
	sessionPub1 := sessionPriv1.PubKey()
	sessionPriv2 := ed25519.GenPrivKey()
	sessionPub2 := sessionPriv2.PubKey()
	
	acc := &gnoland.GnoAccount{
		BaseAccount: std.BaseAccount{
			Address:       addr,
			PubKey:        nil,
			AccountNumber: 1,
			SequenceSum:   0,
		},
		Sessions: []gnoland.GnoSession{
			{
				Name:   "session1",
				PubKey: sessionPub1,
			},
			{
				Name:   "session2",
				PubKey: sessionPub2,
			},
		},
	}
	keeper.setAccount(ctx, acc)

	// Test successful revocation
	t.Run("RevokeSessionSuccess", func(t *testing.T) {
		msg := MsgRevokeSession{
			Signer: addr,
			PubKey: sessionPub1,
		}

		res := handler.Process(ctx, msg)
		assert.True(t, res.IsOK())

		// Verify session was removed
		updatedAcc := keeper.getAccount(ctx, addr).(*gnoland.GnoAccount)
		require.Len(t, updatedAcc.Sessions, 1)
		assert.Equal(t, "session2", updatedAcc.Sessions[0].Name)
	})

	// Test revoking non-existent session
	t.Run("RevokeNonExistentSession", func(t *testing.T) {
		nonExistentPriv := ed25519.GenPrivKey()
		nonExistentPub := nonExistentPriv.PubKey()
		
		msg := MsgRevokeSession{
			Signer: addr,
			PubKey: nonExistentPub,
		}

		res := handler.Process(ctx, msg)
		assert.False(t, res.IsOK())
		assert.Contains(t, res.Err, "session not found")
	})

	// Test revoking from non-existent account
	t.Run("RevokeSessionNonExistentAccount", func(t *testing.T) {
		nonExistentAddr := crypto.Address("nonexistent")
		msg := MsgRevokeSession{
			Signer: nonExistentAddr,
			PubKey: sessionPub2,
		}

		res := handler.Process(ctx, msg)
		assert.False(t, res.IsOK())
		assert.Contains(t, res.Err, "does not exist")
	})
}

func TestSessionHandler_InvalidMessage(t *testing.T) {
	env := setupTestEnv(t)
	ctx := env.ctx
	keeper := env.keeper
	handler := NewSessionHandler(keeper)

	// Test with invalid message type
	invalidMsg := MsgCall{} // Not a session message
	res := handler.Process(ctx, invalidMsg)
	assert.False(t, res.IsOK())
	assert.Contains(t, res.Err, "unrecognized session message type")
}

// Helper to setup test environment
func setupTestEnv(t *testing.T) struct {
	ctx    sdk.Context
	keeper Keeper
} {
	// This is a simplified setup - in real tests you'd use proper test utilities
	t.Helper()
	
	// Mock context and keeper setup would go here
	// For now, returning empty structs as placeholders
	return struct {
		ctx    sdk.Context
		keeper Keeper
	}{
		ctx:    sdk.Context{},
		keeper: Keeper{},
	}
}