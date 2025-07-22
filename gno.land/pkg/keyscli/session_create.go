package keyscli

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys/client"
	"github.com/gnolang/gno/tm2/pkg/std"
)

type SessionCreateCfg struct {
	RootCfg *client.MakeTxCfg

	// Account and key info
	AccountName   string
	SessionName   string  // Name to identify the session locally
	
	// Session properties
	TransferCapacity string
	UnlimitedTransfer bool
	CanManageSessions bool
	CanManagePackages bool
	ValidationOnly bool
	ExpirationDuration string
	RealmsWhitelist string
	InitialSequence uint64
}

func NewSessionCreateCmd(rootCfg *client.MakeTxCfg, io commands.IO) *commands.Command {
	cfg := &SessionCreateCfg{
		RootCfg: rootCfg,
	}

	return commands.NewCommand(
		commands.Metadata{
			Name:       "create",
			ShortUsage: "create [flags]",
			ShortHelp:  "create a new session for an account",
			LongHelp: `Create a new session key for an account with specific permissions.
The session key is generated securely with proper entropy.

Examples:
  # Create a session with 1000ugnot transfer capacity
  gnokey session create -account myaccount -name "daily spending" -transfer-capacity 1000ugnot

  # Create a validation-only session for validators
  gnokey session create -account validator -name "validator session" -validation-only

  # Create a session with realm restrictions
  gnokey session create -account myaccount -name "app session" -realms "r/demo/*,r/myapp/*"`,
		},
		cfg,
		func(ctx context.Context, args []string) error {
			return execSessionCreate(cfg, args, io)
		},
	)
}

func (c *SessionCreateCfg) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.AccountName, "account", "", "account name to create session for")
	fs.StringVar(&c.SessionName, "name", "", "name to identify the session (e.g., 'validator node', 'mobile app')")
	fs.StringVar(&c.TransferCapacity, "transfer-capacity", "", "coin transfer capacity (e.g., 1000ugnot)")
	fs.BoolVar(&c.UnlimitedTransfer, "unlimited-transfer", false, "allow unlimited transfers")
	fs.BoolVar(&c.CanManageSessions, "can-manage-sessions", false, "allow managing other sessions")
	fs.BoolVar(&c.CanManagePackages, "can-manage-packages", false, "allow deploying packages")
	fs.BoolVar(&c.ValidationOnly, "validation-only", false, "restrict to validation operations only")
	fs.StringVar(&c.ExpirationDuration, "expiration", "", "session expiration duration (e.g., 24h, 7d)")
	fs.StringVar(&c.RealmsWhitelist, "realms", "", "comma-separated list of realm patterns to allow")
	fs.Uint64Var(&c.InitialSequence, "sequence", 0, "initial sequence number for the session")
}

func execSessionCreate(cfg *SessionCreateCfg, args []string, io commands.IO) error {
	if cfg.AccountName == "" {
		return fmt.Errorf("account name is required")
	}
	if cfg.SessionName == "" {
		return fmt.Errorf("session name is required")
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

	// Generate a unique session key name based on account and timestamp
	// This ensures we don't collide with existing keys
	timestamp := time.Now().Unix()
	sessionKeyName := fmt.Sprintf("%s_session_%d", cfg.AccountName, timestamp)
	
	// Ensure the generated name doesn't exist (extremely unlikely but safe)
	for i := 0; i < 10; i++ {
		if _, err := kb.GetByName(sessionKeyName); err != nil {
			// Name doesn't exist, we can use it
			break
		}
		// Add a counter to make it unique
		sessionKeyName = fmt.Sprintf("%s_session_%d_%d", cfg.AccountName, timestamp, i)
	}

	// Generate new key for session with secure entropy
	// The empty mnemonic and password mean a new key will be generated
	sessionInfo, err := kb.CreateAccount(sessionKeyName, "", "", "", uint32(0), uint32(0))
	if err != nil {
		return fmt.Errorf("unable to create session key: %w", err)
	}

	// Build the message
	msg := vm.MsgCreateSession{
		Creator:    accountInfo.GetAddress(),
		SessionKey: sessionInfo.GetPubKey(),
	}

	// Set session properties
	if cfg.TransferCapacity != "" {
		coins, err := std.ParseCoins(cfg.TransferCapacity)
		if err != nil {
			return fmt.Errorf("invalid transfer capacity: %w", err)
		}
		msg.TransferCapacity = coins
	}

	msg.UnlimitedTransfer = cfg.UnlimitedTransfer
	msg.CanManageSessions = cfg.CanManageSessions
	msg.CanManagePackages = cfg.CanManagePackages
	msg.ValidationOnly = cfg.ValidationOnly

	if cfg.ExpirationDuration != "" {
		duration, err := time.ParseDuration(cfg.ExpirationDuration)
		if err != nil {
			return fmt.Errorf("invalid expiration duration: %w", err)
		}
		msg.ExpirationTime = time.Now().Add(duration)
	}

	if cfg.RealmsWhitelist != "" {
		// Parse comma-separated realms
		realms := parseCommaSeparated(cfg.RealmsWhitelist)
		msg.RealmsWhitelist = realms
	}

	msg.InitialSequence = cfg.InitialSequence

	// Print session info
	io.Printf("Creating session:\n")
	io.Printf("  Account: %s (%s)\n", cfg.AccountName, accountInfo.GetAddress())
	io.Printf("  Session Name: %s\n", cfg.SessionName)
	io.Printf("  Session Key ID: %s\n", sessionKeyName)
	io.Printf("  Session Address: %s\n", sessionInfo.GetAddress())
	io.Printf("  Session Public Key: %s\n", sessionInfo.GetPubKey())
	
	if cfg.TransferCapacity != "" {
		io.Printf("  Transfer Capacity: %s\n", cfg.TransferCapacity)
	}
	if cfg.UnlimitedTransfer {
		io.Printf("  Unlimited Transfer: yes\n")
	}
	if cfg.CanManageSessions {
		io.Printf("  Can Manage Sessions: yes\n")
	}
	if cfg.CanManagePackages {
		io.Printf("  Can Manage Packages: yes\n")
	}
	if cfg.ValidationOnly {
		io.Printf("  Validation Only: yes\n")
	}
	if cfg.ExpirationDuration != "" {
		io.Printf("  Expires: %s\n", msg.ExpirationTime.Format(time.RFC3339))
	}
	if cfg.RealmsWhitelist != "" {
		io.Printf("  Realms: %s\n", cfg.RealmsWhitelist)
	}
	io.Printf("  Initial Sequence: %d\n", cfg.InitialSequence)
	
	// Parse gas fee
	gasfee, err := std.ParseCoin(cfg.RootCfg.GasFee)
	if err != nil {
		return fmt.Errorf("parsing gas fee coin: %w", err)
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
		err := client.ExecSignAndBroadcast(cfg.RootCfg, args, tx, io)
		if err != nil {
			return err
		}
		io.Println("\nSession created successfully!")
		io.Printf("Session key name: %s\n", sessionKeyName)
		io.Println("Keep this key name safe - you'll need it to use the session.")
	} else {
		io.Println(string(amino.MustMarshalJSON(tx)))
	}

	return nil
}

func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	
	// Simple comma split and trim
	parts := strings.Split(s, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}