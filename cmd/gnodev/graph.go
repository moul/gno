package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gnolang/gno"
	"github.com/gnolang/gno/pkgs/command"
	"github.com/gnolang/gno/pkgs/errors"
	"github.com/gnolang/gno/pkgs/std"
	"github.com/gnolang/gno/pkgs/testutils"
	"github.com/gnolang/gno/tests"
	"go.uber.org/multierr"
)

type graphOptions struct {
	Verbose bool   `flag:"verbose" help:"verbose"`
	RootDir string `flag:"root-dir" help:"clone location of github.com/gnolang/gno (gnodev tries to guess it)"`
}

var DefaultGraphOptions = testOptions{
	Verbose: false,
	RootDir: "",
}

func graphApp(cmd *command.Command, args []string, iopts interface{}) error {
	opts := iopts.(graphOptions)
	if len(args) < 1 {
		cmd.ErrPrintfln("Usage: graph [graph] flags] [packages]")
		return errors.New("invalid args")
	}

	// guess opts.RootDir
	if opts.RootDir == "" {
		opts.RootDir = guessRootDir()
	}

	pkgPaths, err := gnoPackagesFromArgs(args)
	if err != nil {
		return fmt.Errorf("list packages from args: %w", err)
	}

	// FIXME: get the source (local, gno.mod replace, fetched from the chain)
	return nil
}
