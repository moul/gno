package keyscli

import (
	"context"
	"flag"
	"fmt"

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys/client"
	"github.com/gnolang/gno/tm2/pkg/std"
)

type SessionRevokeCfg struct {
	RootCfg *client.MakeTxCfg

	AccountName    string
	SessionKeyName string  // The generated key name from keybase
}

func NewSessionRevokeCmd(rootCfg *client.MakeTxCfg, io commands.IO) *commands.Command {
	cfg := &SessionRevokeCfg{
		RootCfg: rootCfg,
	}

	return commands.NewCommand(
		commands.Metadata{
			Name:       "revoke",
			ShortUsage: "revoke [flags]",
			ShortHelp:  "revoke a session from an account",
			LongHelp:   `Revoke a session key from an account, removing its permissions.`,
		},
		cfg,
		func(ctx context.Context, args []string) error {
			return execSessionRevoke(cfg, args, io)
		},
	)
}

func (c *SessionRevokeCfg) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.AccountName, "account", "", "account name to revoke session from")
	fs.StringVar(&c.SessionKeyName, "session-key-name", "", "keybase name of the session key to revoke (from session create output)")
}

func execSessionRevoke(cfg *SessionRevokeCfg, args []string, io commands.IO) error {
	if cfg.AccountName == "" {
		return fmt.Errorf("account name is required")
	}
	if cfg.SessionKeyName == "" {
		return fmt.Errorf("session key name is required")
	}

	// Check gas parameters
	if cfg.RootCfg.GasWanted == 0 {
		return fmt.Errorf("gas-wanted not specified")
	}
	if cfg.RootCfg.GasFee == "" {
		return fmt.Errorf("gas-fee not specified")
	}

	// Open keybase
	kb, err := keys.NewKeyBaseFromDir(cfg.RootCfg.RootCfg.Home)
	if err != nil {
		return fmt.Errorf("unable to open keybase: %w", err)
	}

	// Get account info
	accountInfo, err := kb.GetByName(cfg.AccountName)
	if err != nil {
		return fmt.Errorf("unable to get account info: %w", err)
	}

	// Get session key info
	sessionInfo, err := kb.GetByName(cfg.SessionKeyName)
	if err != nil {
		return fmt.Errorf("unable to get session key info: %w", err)
	}

	io.Printf("Revoking session:\n")
	io.Printf("  Account: %s (%s)\n", cfg.AccountName, accountInfo.GetAddress())
	io.Printf("  Session Key Name: %s\n", cfg.SessionKeyName)
	io.Printf("  Session Address: %s\n", sessionInfo.GetAddress())

	// Parse gas fee
	gasfee, err := std.ParseCoin(cfg.RootCfg.GasFee)
	if err != nil {
		return fmt.Errorf("parsing gas fee coin: %w", err)
	}

	// Create message
	msg := vm.MsgRevokeSession{
		Creator:    accountInfo.GetAddress(),
		SessionKey: sessionInfo.GetPubKey(),
	}

	// Construct transaction
	tx := std.Tx{
		Msgs:       []std.Msg{msg},
		Fee:        std.NewFee(cfg.RootCfg.GasWanted, gasfee),
		Signatures: nil,
		Memo:       cfg.RootCfg.Memo,
	}

	// Sign and broadcast the transaction
	if cfg.RootCfg.Broadcast {
		err := client.ExecSignAndBroadcast(cfg.RootCfg, []string{cfg.AccountName}, tx, io)
		if err != nil {
			return err
		}
		io.Println("\nSession revoked successfully!")
		io.Println("After confirming revocation, you may also want to delete the key from keybase:")
		io.Printf("  gnokey delete %s\n", cfg.SessionKeyName)
	} else {
		io.Println(string(amino.MustMarshalJSON(tx)))
	}

	return nil
}