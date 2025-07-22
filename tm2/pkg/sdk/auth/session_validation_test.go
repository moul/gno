package auth

import (
	"testing"
	"time"

	"github.com/gnolang/gno/gno.land/pkg/gnoland"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/crypto/ed25519"
	"github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/sdk/bank"
	"github.com/gnolang/gno/tm2/pkg/std"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateSessionForTx(t *testing.T) {
	// Setup test keys
	masterPriv := ed25519.GenPrivKey()
	masterPub := masterPriv.PubKey()
	masterAddr := masterPub.Address()

	sessionPriv := ed25519.GenPrivKey()
	sessionPub := sessionPriv.PubKey()

	// Create test context with block time
	ctx := sdk.Context{}.WithBlockTime(time.Now())

	t.Run("MasterKeyAlwaysValid", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
		}

		// Any transaction should be valid with master key
		tx := std.Tx{
			Msgs: []std.Msg{
				bank.MsgSend{
					FromAddress: masterAddr,
					ToAddress:   crypto.Address("recipient"),
					Amount:      std.Coins{std.NewCoin("ugnot", 1000000)},
				},
			},
		}

		res := ValidateSessionForTx(ctx, acc, masterPub, tx)
		assert.True(t, res.IsOK())
	})

	t.Run("SessionWithSufficientCapacity", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
			Sessions: []gnoland.GnoSession{
				{
					Name:             "test-session",
					PubKey:           sessionPub,
					TransferCapacity: 1000000, // 1 GNOT
				},
			},
		}

		// Transaction within capacity
		tx := std.Tx{
			Msgs: []std.Msg{
				bank.MsgSend{
					FromAddress: masterAddr,
					ToAddress:   crypto.Address("recipient"),
					Amount:      std.Coins{std.NewCoin("ugnot", 500000)}, // 0.5 GNOT
				},
			},
		}

		res := ValidateSessionForTx(ctx, acc, sessionPub, tx)
		assert.True(t, res.IsOK())
	})

	t.Run("SessionExceedsCapacity", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
			Sessions: []gnoland.GnoSession{
				{
					Name:             "limited-session",
					PubKey:           sessionPub,
					TransferCapacity: 100000, // 0.1 GNOT
				},
			},
		}

		// Transaction exceeding capacity
		tx := std.Tx{
			Msgs: []std.Msg{
				bank.MsgSend{
					FromAddress: masterAddr,
					ToAddress:   crypto.Address("recipient"),
					Amount:      std.Coins{std.NewCoin("ugnot", 200000)}, // 0.2 GNOT
				},
			},
		}

		res := ValidateSessionForTx(ctx, acc, sessionPub, tx)
		assert.False(t, res.IsOK())
		assert.Contains(t, res.Err, "session transfer capacity exceeded")
	})

	t.Run("ExpiredSession", func(t *testing.T) {
		expireTime := ctx.BlockTime().Add(-1 * time.Hour) // Expired 1 hour ago
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
			Sessions: []gnoland.GnoSession{
				{
					Name:     "expired-session",
					PubKey:   sessionPub,
					ExpireAt: &expireTime,
				},
			},
		}

		tx := std.Tx{
			Msgs: []std.Msg{
				bank.MsgSend{
					FromAddress: masterAddr,
					ToAddress:   crypto.Address("recipient"),
					Amount:      std.Coins{std.NewCoin("ugnot", 1000)},
				},
			},
		}

		res := ValidateSessionForTx(ctx, acc, sessionPub, tx)
		assert.False(t, res.IsOK())
		assert.Contains(t, res.Err, "session has expired")
	})

	t.Run("SessionNotFound", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
			Sessions: []gnoland.GnoSession{}, // No sessions
		}

		tx := std.Tx{
			Msgs: []std.Msg{
				bank.MsgSend{
					FromAddress: masterAddr,
					ToAddress:   crypto.Address("recipient"),
					Amount:      std.Coins{std.NewCoin("ugnot", 1000)},
				},
			},
		}

		unknownSessionPriv := ed25519.GenPrivKey()
		unknownSessionPub := unknownSessionPriv.PubKey()

		res := ValidateSessionForTx(ctx, acc, unknownSessionPub, tx)
		assert.False(t, res.IsOK())
		assert.Contains(t, res.Err, "session not found")
	})

	t.Run("MultiSendTransaction", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
			Sessions: []gnoland.GnoSession{
				{
					Name:             "multisend-session",
					PubKey:           sessionPub,
					TransferCapacity: 1000000, // 1 GNOT
				},
			},
		}

		// MultiSend with total within capacity
		tx := std.Tx{
			Msgs: []std.Msg{
				bank.MsgMultiSend{
					Inputs: []bank.Input{
						{
							Address: masterAddr,
							Coins:   std.Coins{std.NewCoin("ugnot", 600000)},
						},
					},
					Outputs: []bank.Output{
						{
							Address: crypto.Address("recipient1"),
							Coins:   std.Coins{std.NewCoin("ugnot", 300000)},
						},
						{
							Address: crypto.Address("recipient2"),
							Coins:   std.Coins{std.NewCoin("ugnot", 300000)},
						},
					},
				},
			},
		}

		res := ValidateSessionForTx(ctx, acc, sessionPub, tx)
		assert.True(t, res.IsOK())
	})
}

func TestConsumeSessionTransferCapacity(t *testing.T) {
	masterPriv := ed25519.GenPrivKey()
	masterPub := masterPriv.PubKey()
	masterAddr := masterPub.Address()

	sessionPriv := ed25519.GenPrivKey()
	sessionPub := sessionPriv.PubKey()

	ctx := sdk.Context{}

	t.Run("ConsumeCapacitySuccess", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
			Sessions: []gnoland.GnoSession{
				{
					Name:             "test-session",
					PubKey:           sessionPub,
					TransferCapacity: 1000000,
				},
			},
		}

		err := ConsumeSessionTransferCapacity(ctx, acc, sessionPub, 400000)
		require.NoError(t, err)

		// Verify capacity was reduced
		session, err := acc.GetSession(sessionPub)
		require.NoError(t, err)
		assert.Equal(t, int64(600000), session.TransferCapacity)
	})

	t.Run("ConsumeCapacityMasterKey", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
		}

		// Master key should not consume capacity
		err := ConsumeSessionTransferCapacity(ctx, acc, masterPub, 1000000)
		assert.NoError(t, err)
	})

	t.Run("ConsumeCapacityInsufficientFunds", func(t *testing.T) {
		acc := &gnoland.GnoAccount{
			BaseAccount: std.BaseAccount{
				Address: masterAddr,
				Keys: []std.AccountKey{
					&gnoland.BaseAccountKey{
						Address:  masterAddr,
						PubKey:   masterPub,
						Sequence: 0,
					},
				},
			},
			Sessions: []gnoland.GnoSession{
				{
					Name:             "limited-session",
					PubKey:           sessionPub,
					TransferCapacity: 100000,
				},
			},
		}

		err := ConsumeSessionTransferCapacity(ctx, acc, sessionPub, 200000)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient transfer capacity")
	})
}

func TestIsSessionExpired(t *testing.T) {
	now := time.Now()

	t.Run("NotExpired", func(t *testing.T) {
		futureTime := now.Add(1 * time.Hour)
		session := &gnoland.GnoSession{
			ExpireAt: &futureTime,
		}
		assert.False(t, IsSessionExpired(session, now))
	})

	t.Run("Expired", func(t *testing.T) {
		pastTime := now.Add(-1 * time.Hour)
		session := &gnoland.GnoSession{
			ExpireAt: &pastTime,
		}
		assert.True(t, IsSessionExpired(session, now))
	})

	t.Run("NoExpiration", func(t *testing.T) {
		session := &gnoland.GnoSession{
			ExpireAt: nil,
		}
		assert.False(t, IsSessionExpired(session, now))
	})
}