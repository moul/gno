package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/tests"
	"github.com/gnolang/gno/tm2/pkg/commands"
	"golang.org/x/term"
)

type replCfg struct {
	verbose        bool
	rootDir        string
	initialImports string
	initialCommand string
	skipUsage      bool
}

func newReplCmd() *commands.Command {
	cfg := &replCfg{}

	return commands.NewCommand(
		commands.Metadata{
			Name:       "repl",
			ShortUsage: "repl [flags]",
			ShortHelp:  "Starts a GnoVM REPL",
		},
		cfg,
		func(_ context.Context, args []string) error {
			return execRepl(cfg, args)
		},
	)
}

func (c *replCfg) RegisterFlags(fs *flag.FlagSet) {
	fs.BoolVar(
		&c.verbose,
		"verbose",
		false,
		"verbose output when running",
	)

	fs.StringVar(
		&c.rootDir,
		"root-dir",
		"",
		"clone location of github.com/gnolang/gno (gnodev tries to guess it)",
	)

	fs.StringVar(
		&c.initialImports,
		"imports",
		"gno.land/p/demo/avl,gno.land/p/demo/ufmt",
		"initial imports, separated by a comma",
	)

	fs.StringVar(
		&c.initialCommand,
		"command",
		"",
		"initial command to run",
	)

	fs.BoolVar(
		&c.skipUsage,
		"skip-usage",
		false,
		"do not print usage",
	)
}

func execRepl(cfg *replCfg, args []string) error {
	if len(args) > 0 {
		return flag.ErrHelp
	}

	if cfg.rootDir == "" {
		cfg.rootDir = guessRootDir()
	}

	if !cfg.skipUsage {
		fmt.Fprint(os.Stderr, `// Usage:
//   gno:1> /import "gno.land/p/demo/avl"     // import the p/demo/avl package
//   gno:2> /func a() string { return "a" }   // declare a new function named a
//   gno:3> /src                              // print current generated source
//   gno:4> println(a())                      // print the result of calling a()
//   gno:5> /exit
`)
	}

	return runRepl(cfg)
}

type repl struct {
	// repl state
	imports   []string
	funcs     []string
	lastInput string
	i         int
	// TODO: support setting global vars
	// TODO: switch to state machine, and support rollback of anything

	machine *gno.Machine
}

func (r *repl) handleInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	r.i++
	funcName := fmt.Sprintf("repl_%d", r.i)
	// FIXME: support ";" as line separator?
	// FIXME: support multiline when unclosed parenthesis, etc

	imports := strings.Join(r.imports, "\n")
	funcs := strings.Join(r.funcs, "\n")
	src := "package test\n" + imports + "\n" + funcs + "\nfunc " + funcName + "() {\nINPUT\n}"

	fields := strings.Fields(input)
	command := fields[0]
	switch {
	case command == "/import":
		imp := fields[1]
		r.imports = append(r.imports, `import "`+imp+`"`)
		// TODO: check if valid, else rollback
		return nil
	case command == "/func":
		r.funcs = append(r.funcs, input[1:])
		// TODO: check if valid, else rollback
		return nil
	case command == "/src":
		// TODO: use go/format for pretty print
		src = strings.ReplaceAll(src, "INPUT", r.lastInput)
		println(src)
		return nil
	case command == "/exit":
		os.Exit(0) // return special err?
	case strings.HasPrefix(command, "/"):
		println("unsupported command")
		return nil
	default:
		// not a command, probably code to run
	}

	r.lastInput = input
	src = strings.ReplaceAll(src, "INPUT", input)
	n := gno.MustParseFile(funcName+".gno", src)
	// TODO: run fmt check + linter
	r.machine.RunFiles(n)
	// TODO: smart recover system
	r.machine.RunStatement(gno.S(gno.Call(gno.X(funcName))))
	// TODO: if output is empty, consider that it's a persisted variable?
	return nil
}

func runRepl(cfg *replCfg) error {
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	// init repl state
	r := repl{
		imports:   make([]string, 0),
		funcs:     make([]string, 0),
		lastInput: "// your code will be here", // initial value, to make it easier to identify with '/src'
	}
	for _, imp := range strings.Split(cfg.initialImports, ",") {
		if strings.TrimSpace(imp) == "" {
			continue
		}
		r.imports = append(r.imports, `import "`+imp+`"`)
	}
	if cfg.initialCommand != "" {
		// TODO: implement
		panic("not implemented")
	}
	testStore := tests.TestStore(cfg.rootDir, "", stdin, stdout, stderr, tests.ImportModeStdlibsOnly)
	if cfg.verbose {
		testStore.SetLogStoreOps(true)
	}
	r.machine = gno.NewMachineWithOptions(gno.MachineOptions{
		PkgPath: "test",
		Output:  stdout,
		Store:   testStore,
	})
	defer r.machine.Release()

	// main loop
	isTerm := term.IsTerminal(int(stdin.Fd()))

	if isTerm {
		rw := struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stderr}
		t := term.NewTerminal(rw, "")
		for i := 1; ; i++ {
			// prompt and parse
			t.SetPrompt(fmt.Sprintf("gno:%d> ", i))
			oldState, err := term.MakeRaw(0)
			if err != nil {
				return fmt.Errorf("make term raw: %w", err)
			}
			input, err := t.ReadLine()
			if err != nil {
				term.Restore(0, oldState)
				if errors.Is(err, io.EOF) {
					return nil
				}
				return fmt.Errorf("term error: %w", err)
			}
			term.Restore(0, oldState)

			err = r.handleInput(input)
			if err != nil {
				return fmt.Errorf("handle repl input: %w", err)
			}
		}
	} else { // !isTerm
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			input := scanner.Text()
			err := r.handleInput(input)
			if err != nil {
				return fmt.Errorf("handle repl input: %w", err)
			}
		}
		err := scanner.Err()
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
	}
	return nil
}
