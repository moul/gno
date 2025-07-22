package keyscli

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys/client"
)

type SessionListCfg struct {
	RootCfg *client.BaseCfg

	AccountName string
}

func NewSessionListCmd(rootCfg *client.BaseCfg, io commands.IO) *commands.Command {
	cfg := &SessionListCfg{
		RootCfg: rootCfg,
	}

	return commands.NewCommand(
		commands.Metadata{
			Name:       "list",
			ShortUsage: "list [flags]",
			ShortHelp:  "list sessions for an account",
			LongHelp:   `List all active sessions for a given account, showing their permissions and status.
This shows sessions created by this keybase, identified by the naming pattern.`,
		},
		cfg,
		func(ctx context.Context, args []string) error {
			return execSessionList(cfg, args, io)
		},
	)
}

func (c *SessionListCfg) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.AccountName, "account", "", "account name to list sessions for")
}

func execSessionList(cfg *SessionListCfg, args []string, io commands.IO) error {
	if cfg.AccountName == "" {
		return fmt.Errorf("account name is required")
	}

	// Open keybase
	kb, err := keys.NewKeyBaseFromDir(cfg.RootCfg.Home)
	if err != nil {
		return fmt.Errorf("unable to open keybase: %w", err)
	}

	// Get account info
	accountInfo, err := kb.GetByName(cfg.AccountName)
	if err != nil {
		return fmt.Errorf("unable to get account info: %w", err)
	}

	io.Printf("Account: %s (%s)\n\n", cfg.AccountName, accountInfo.GetAddress())
	io.Println("Local Session Keys:")
	
	// List all keys in the keybase and find session keys for this account
	allKeys, err := kb.List()
	if err != nil {
		return fmt.Errorf("unable to list keys: %w", err)
	}
	
	sessionPrefix := cfg.AccountName + "_session_"
	sessionCount := 0
	
	for _, keyInfo := range allKeys {
		keyName := keyInfo.GetName()
		if strings.HasPrefix(keyName, sessionPrefix) {
			sessionCount++
			io.Printf("%d. Key Name: %s\n", sessionCount, keyName)
			io.Printf("   Address: %s\n", keyInfo.GetAddress())
			io.Printf("   Public Key: %s\n", keyInfo.GetPubKey())
			io.Println()
		}
	}
	
	if sessionCount == 0 {
		io.Println("No session keys found for this account.")
		io.Println("Create one with: gnokey session create -account", cfg.AccountName, "-name <session-name>")
	}
	
	// TODO: Query blockchain for account sessions to show their permissions and status
	io.Println("NOTE: On-chain session status querying not yet implemented")
	io.Println("This will query the blockchain for active sessions and their permissions")

	return nil
}