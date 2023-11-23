package gnoclient

import (
	"fmt"
	"testing"

	"github.com/gnolang/gno/gno.land/pkg/gnoland"
	"github.com/gnolang/gno/gno.land/pkg/integration"
	rpcclient "github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/log"
	"github.com/jaekwon/testify/require"
)

func newInMemorySigner(t *testing.T, chainid string) *SignerFromKeybase {
	t.Helper()

	mmeonic := integration.DefaultAccount_Seed
	name := integration.DefaultAccount_Name

	kb := keys.NewInMemory()
	_, err := kb.CreateAccount(name, mmeonic, "", "", uint32(0), uint32(0))
	require.NoError(t, err)

	return &SignerFromKeybase{
		Keybase:  kb,      // Stores keys in memory or on disk
		Account:  name,    // Account name or bech32 format
		Password: "",      // Password for encryption
		ChainID:  chainid, // Chain ID for transaction signing
	}
}

func TestClient_Request(t *testing.T) {
	config, _ := integration.TestingNodeConfig(t, gnoland.MustGuessGnoRootDir())
	node, remoteAddr := integration.TestingInMemoryNode(t, log.NewNopLogger(), config)
	defer node.Stop()

	signer := newInMemorySigner(t, config.TMConfig.ChainID())

	client := Client{
		Signer:    signer,
		RPCClient: rpcclient.NewHTTP(remoteAddr, "/websocket"),
	}

	data, res, err := client.Render("gno.land/r/demo/boards", "")
	require.NoError(t, err)
	fmt.Println("data", data)
	fmt.Println("res", res)

	// TODO: xxx
}
