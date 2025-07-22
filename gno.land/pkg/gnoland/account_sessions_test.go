package gnoland

import (
	"testing"
	"time"

	"github.com/gnolang/gno/tm2/pkg/crypto/secp256k1"
	"github.com/gnolang/gno/tm2/pkg/std"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountSessions(t *testing.T) {
	// Create a new account with a master key
	masterPrivKey := secp256k1.GenPrivKey()
	masterPubKey := masterPrivKey.PubKey()
	accountAddr := masterPubKey.Address()

	account := NewGnoAccountWithMasterKey(accountAddr, masterPubKey)

	// Test 1: Create a session
	sessionPrivKey := secp256k1.GenPrivKey()
	sessionPubKey := sessionPrivKey.PubKey()

	session, err := account.CreateSession(sessionPubKey)
	require.NoError(t, err, "Should create session successfully")
	assert.NotNil(t, session, "Session should not be nil")
	assert.True(t, session.MatchesPubKey(sessionPubKey), "Session should match pubkey")

	// Test 2: List sessions
	sessions := account.GetSessions()
	assert.Len(t, sessions, 1, "Should have one session")
	assert.True(t, sessions[0].MatchesPubKey(sessionPubKey), "Listed session should match")

	// Test 3: Get specific session
	retrievedSession, err := account.GetSession(sessionPubKey)
	require.NoError(t, err, "Should get session successfully")
	assert.True(t, retrievedSession.MatchesPubKey(sessionPubKey), "Retrieved session should match")

	// Test 4: Configure session permissions
	retrievedSession.SetCoinsTransferCapacity(std.Coins{std.NewCoin("ugnot", 1000)})
	retrievedSession.SetValidationOnly()
	retrievedSession.SetExpirationTime(time.Now().Add(24 * time.Hour))
	retrievedSession.SetRealmsWhitelist([]string{"r/demo/*", "r/myapp/*"})

	// Test 5: Verify session permissions
	assert.True(t, retrievedSession.IsValidationOnly(), "Should be validation only")
	assert.False(t, retrievedSession.CanManageSessions(), "Should not be able to manage sessions")
	
	// Test transfer capacity
	smallAmount := std.Coins{std.NewCoin("ugnot", 100)}
	largeAmount := std.Coins{std.NewCoin("ugnot", 2000)}
	assert.True(t, retrievedSession.CanTransferAmount(smallAmount), "Should allow small transfer")
	assert.False(t, retrievedSession.CanTransferAmount(largeAmount), "Should deny large transfer")

	// Test realm access
	assert.True(t, retrievedSession.HasRealmAccess("r/demo/board"), "Should have access to r/demo/board")
	assert.True(t, retrievedSession.HasRealmAccess("r/myapp/test"), "Should have access to r/myapp/test")
	assert.False(t, retrievedSession.HasRealmAccess("r/other/realm"), "Should not have access to r/other/realm")

	// Test 6: Consume transfer capacity
	err = retrievedSession.ConsumeTransferCapacity(smallAmount)
	require.NoError(t, err, "Should consume capacity successfully")
	
	remainingCapacity := retrievedSession.CoinsTransferCapacity
	expectedRemaining := std.Coins{std.NewCoin("ugnot", 900)}
	assert.True(t, remainingCapacity.IsEqual(expectedRemaining), "Remaining capacity should be correct")

	// Test 7: Revoke session
	err = account.RevokeSession(sessionPubKey)
	require.NoError(t, err, "Should revoke session successfully")

	sessions = account.GetSessions()
	assert.Len(t, sessions, 0, "Should have no sessions after revocation")

	// Test 8: Try to get revoked session
	_, err = account.GetSession(sessionPubKey)
	assert.Error(t, err, "Should error when getting revoked session")
}

func TestSessionExpiration(t *testing.T) {
	masterPrivKey := secp256k1.GenPrivKey()
	masterPubKey := masterPrivKey.PubKey()
	accountAddr := masterPubKey.Address()

	account := NewGnoAccountWithMasterKey(accountAddr, masterPubKey)

	// Create an expired session
	sessionPubKey := secp256k1.GenPrivKey().PubKey()
	session, err := account.CreateSession(sessionPubKey)
	require.NoError(t, err)

	// Set expiration in the past
	session.SetExpirationTime(time.Now().Add(-1 * time.Hour))

	// Session should be expired
	assert.True(t, session.IsExpired(), "Session should be expired")

	// Getting expired session should fail
	_, err = account.GetSession(sessionPubKey)
	require.Error(t, err, "Should error when getting expired session")
	assert.Contains(t, err.Error(), "expired", "Error should mention expiration")

	// Garbage collection should remove it
	removed := account.gc()
	assert.Equal(t, 1, removed, "Should have removed one expired session")

	sessions := account.GetSessions()
	assert.Len(t, sessions, 0, "Should have no sessions after gc")
}

func TestSessionWithCustomSequence(t *testing.T) {
	masterPrivKey := secp256k1.GenPrivKey()
	masterPubKey := masterPrivKey.PubKey()
	accountAddr := masterPubKey.Address()

	// Create session with custom initial sequence
	sessionPubKey := secp256k1.GenPrivKey().PubKey()
	customSequence := uint64(1000)
	
	session := NewGnoSessionWithSequence(accountAddr, sessionPubKey, customSequence)
	
	assert.Equal(t, customSequence, session.GetSequence(), "Session should have custom sequence")
	
	// Test sequence increment
	err := session.SetSequence(customSequence + 1)
	require.NoError(t, err)
	assert.Equal(t, customSequence+1, session.GetSequence(), "Sequence should be incremented")
}

func TestMaxSessionsLimit(t *testing.T) {
	masterPrivKey := secp256k1.GenPrivKey()
	masterPubKey := masterPrivKey.PubKey()
	accountAddr := masterPubKey.Address()

	account := NewGnoAccountWithMasterKey(accountAddr, masterPubKey)

	// Create maximum number of sessions
	for i := 0; i < MaxSessionsPerAccount; i++ {
		sessionPubKey := secp256k1.GenPrivKey().PubKey()
		_, err := account.CreateSession(sessionPubKey)
		require.NoError(t, err, "Should create session %d", i)
	}

	// Try to create one more
	extraPubKey := secp256k1.GenPrivKey().PubKey()
	_, err := account.CreateSession(extraPubKey)
	assert.Error(t, err, "Should error when exceeding max sessions")
	assert.Contains(t, err.Error(), "maximum", "Error should mention maximum")
}

func TestDuplicateSessionKey(t *testing.T) {
	masterPrivKey := secp256k1.GenPrivKey()
	masterPubKey := masterPrivKey.PubKey()
	accountAddr := masterPubKey.Address()

	account := NewGnoAccountWithMasterKey(accountAddr, masterPubKey)

	// Create a session
	sessionPubKey := secp256k1.GenPrivKey().PubKey()
	_, err := account.CreateSession(sessionPubKey)
	require.NoError(t, err)

	// Try to create another session with same key
	_, err = account.CreateSession(sessionPubKey)
	assert.Error(t, err, "Should error when creating duplicate session")
	assert.Contains(t, err.Error(), "already exists", "Error should mention duplicate")
}