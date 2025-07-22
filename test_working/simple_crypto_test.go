package test_working

import (
	"testing"

	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCryptoBasics(t *testing.T) {
	// Test basic crypto functionality without amino initialization
	addr := crypto.Bech32Address("g1234567890abcdef")
	assert.NotEmpty(t, addr.String())
	assert.Equal(t, "g1234567890abcdef", addr.String())
}