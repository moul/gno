package keyscli

import (
	"context"
	"errors"

	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys/client"
)

// NewSessionCmd creates the session management parent command
func NewSessionCmd(rootCfg *client.BaseCfg, io commands.IO) *commands.Command {
	// Create a MakeTxCfg for subcommands that need it
	makeTxCfg := &client.MakeTxCfg{
		RootCfg: rootCfg,
	}

	cmd := commands.NewCommand(
		commands.Metadata{
			Name:       "session",
			ShortUsage: "session <subcommand> [flags]",
			ShortHelp:  "manage account sessions",
			LongHelp: `Account sessions allow delegating limited permissions to different keys.
Sessions can have:
- Transfer capacity limits
- Realm access restrictions
- Expiration times
- Specific permission flags`,
		},
		makeTxCfg,
		func(ctx context.Context, args []string) error {
			return errors.New("subcommand required")
		},
	)

	// Add subcommands
	cmd.AddSubCommands(
		NewSessionCreateCmd(makeTxCfg, io),
		NewSessionListCmd(rootCfg, io),
		NewSessionRevokeCmd(makeTxCfg, io),
	)

	return cmd
}