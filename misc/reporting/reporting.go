package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"os"
)

type Opts struct {
	changelog    bool
	backlog      bool
	curation     bool
	tips         bool
	From         string
	twitterToken string
	githubToken  string
	help         bool
}

var opts = NewOpts()

func main() {
	err := runMain(os.Args[1:])
	if err == flag.ErrHelp {
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

func runMain(args []string) error {
	var root *ffcli.Command
	{
		globalFlags := flag.NewFlagSet("gno-reporting", flag.ExitOnError)
		globalFlags.BoolVar(&opts.changelog, "changelog", opts.changelog, "generate changelog")
		globalFlags.BoolVar(&opts.backlog, "backlog", opts.backlog, "generate backlog")
		globalFlags.BoolVar(&opts.curation, "curation", opts.curation, "generate curation")
		globalFlags.BoolVar(&opts.tips, "tips", opts.tips, "generate tips")
		globalFlags.StringVar(&opts.From, "from", opts.From, "from date")
		globalFlags.StringVar(&opts.twitterToken, "twitter-token", opts.twitterToken, "twitter token")
		globalFlags.StringVar(&opts.githubToken, "github-token", opts.githubToken, "github token")
		globalFlags.BoolVar(&opts.help, "help", false, "show help")
		root = &ffcli.Command{
			ShortUsage: "reporting [flags]",
			FlagSet:    globalFlags,
			Exec: func(ctx context.Context, args []string) error {
				if opts.help {
					return flag.ErrHelp
				}
				fmt.Println("hello world")
				return nil
			},
		}
	}
	return root.ParseAndRun(context.Background(), args)
}

func NewOpts() Opts {
	return Opts{
		changelog:    true,
		backlog:      true,
		curation:     true,
		tips:         true,
		From:         "",
		twitterToken: "",
		githubToken:  "",
		help:         false,
	}
}
