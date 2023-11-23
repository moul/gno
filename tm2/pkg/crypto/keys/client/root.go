// Dedicated to my love, Lexi.
package client

import (
	"flag"

	"github.com/gnolang/gno/tm2/pkg/commands"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/fftoml"
)

const (
	mnemonicEntropySize = 256
)

type baseCfg struct {
	BaseOptions
}

func NewRootCmd(io commands.IO) *commands.Command {
	cfg := &baseCfg{}

	cmd := commands.NewCommand(
		commands.Metadata{
			ShortUsage: "<subcommand> [flags] [<arg>...]",
			LongHelp:   "Manages private keys for the node",
			Options: []ff.Option{
				ff.WithConfigFileFlag("config"),
				ff.WithConfigFileParser(fftoml.Parser),
			},
		},
		cfg,
		commands.HelpExec,
	)

	cmd.AddSubCommands(
		newAddCmd(cfg, io),
		newDeleteCmd(cfg, io),
		newGenerateCmd(cfg, io),
		newExportCmd(cfg, io),
		newImportCmd(cfg, io),
		newListCmd(cfg, io),
		newSignCmd(cfg, io),
		newVerifyCmd(cfg, io),
		newQueryCmd(cfg, io),
		newBroadcastCmd(cfg, io),
		newMakeTxCmd(cfg, io),
	)

	return cmd
}

func (c *baseCfg) RegisterFlags(fs *flag.FlagSet) {
	// Base options
	fs.StringVar(
		&c.Home,
		"home",
		DefaultBaseOptions.Home,
		"home directory",
	)

	fs.StringVar(
		&c.Remote,
		"remote",
		DefaultBaseOptions.Remote,
		"remote node URL",
	)

	fs.BoolVar(
		&c.Quiet,
		"quiet",
		DefaultBaseOptions.Quiet,
		"suppress output during execution",
	)

	fs.BoolVar(
		&c.InsecurePasswordStdin,
		"insecure-password-stdin",
		DefaultBaseOptions.Quiet,
		"WARNING! take password from stdin",
	)

	fs.StringVar(
		&c.Config,
		"config",
		DefaultBaseOptions.Config,
		"config file (optional)",
	)
}
