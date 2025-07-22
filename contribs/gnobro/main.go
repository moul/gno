package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/signal"
	gopath "path"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	charmlog "github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"golang.org/x/sync/errgroup"

	"github.com/gnolang/gno/contribs/gnobro/pkg/browser"
	"github.com/gnolang/gno/contribs/gnodev/pkg/events"
	"github.com/gnolang/gno/gno.land/pkg/gnoclient"
	"github.com/gnolang/gno/gno.land/pkg/integration"
	"github.com/gnolang/gno/gnovm/pkg/gnoenv"
	"github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
)

const gnoPrefix = "gno.land"

// Common configuration shared between remote and local modes
type commonCfg struct {
	readonly       bool
	defaultAccount string
	defaultRealm   string
	sshListener    string
	sshHostKeyPath string
	banner         bool
	jsonlog        bool
}

// Remote configuration for connecting to standard Gno.land chains
type remoteCfg struct {
	commonCfg
	remote  string
	chainID string
}

// Local configuration for connecting to gnodev
type localCfg struct {
	commonCfg
	endpoint string
}

var defaultCommonOptions = commonCfg{
	defaultRealm:   "gno.land/r/gnoland/home",
	sshHostKeyPath: ".ssh/id_ed25519",
}

var defaultRemoteOptions = remoteCfg{
	commonCfg: defaultCommonOptions,
	remote:    "https://rpc.gno.land:443",
	chainID:   "portal-loop",
}

var defaultLocalOptions = localCfg{
	commonCfg:      defaultCommonOptions,
	endpoint: "ws://127.0.0.1:8888",
}

func main() {
	stdio := commands.NewDefaultIO()
	
	// Main command (defaults to remote mode)
	remoteCfg := &remoteCfg{}
	cmd := commands.NewCommand(
		commands.Metadata{
			Name:       "gnobro",
			ShortUsage: "gnobro [flags] [realm_path]",
			ShortHelp:  "Gno Browser - browse Gno.land realms",
			LongHelp: `Gnobro is a terminal user interface (TUI) for browsing Gno.land realms.

By default, gnobro connects to the remote Gno.land chain.

Examples:
  gnobro                                    # Browse default realm
  gnobro gno.land/r/demo/boards            # Browse specific realm
  gnobro local                             # Connect to local gnodev
  gnobro local gno.land/r/demo/boards      # Browse realm on local gnodev`,
		},
		remoteCfg,
		func(ctx context.Context, args []string) error {
			return execRemote(ctx, remoteCfg, args, stdio)
		})

	// Add local subcommand
	localCfg := &localCfg{}
	localCmd := commands.NewCommand(
		commands.Metadata{
			Name:       "local",
			ShortUsage: "gnobro local [flags] [realm_path]",
			ShortHelp:  "Connect to local gnodev instance",
			LongHelp: `Connect to a local gnodev instance for development.

This mode provides:
- Hot reload on code changes
- Real-time event monitoring
- Transaction result notifications
- Development-focused UI

The connection is established through gnodev's WebSocket endpoint.

Examples:
  gnobro local                             # Connect to default gnodev
  gnobro local --endpoint ws://localhost:9999  # Custom endpoint`,
		},
		localCfg,
		func(ctx context.Context, args []string) error {
			return execLocal(ctx, localCfg, args, stdio)
		})
	
	cmd.AddSubCommands(localCmd)
	cmd.Execute(context.Background(), os.Args[1:])
}

// RegisterFlags for remoteCfg
func (c *remoteCfg) RegisterFlags(fs *flag.FlagSet) {
	// Common flags
	c.commonCfg.RegisterFlags(fs)
	
	// Remote-specific flags
	fs.StringVar(
		&c.remote,
		"remote",
		defaultRemoteOptions.remote,
		"remote gno.land RPC URL",
	)

	fs.StringVar(
		&c.chainID,
		"chainid",
		defaultRemoteOptions.chainID,
		"chain ID",
	)
}

// RegisterFlags for localCfg
func (c *localCfg) RegisterFlags(fs *flag.FlagSet) {
	// Common flags
	c.commonCfg.RegisterFlags(fs)
	
	// Local-specific flags
	fs.StringVar(
		&c.endpoint,
		"endpoint",
		defaultLocalOptions.endpoint,
		"WebSocket endpoint for gnodev events",
	)
}

// RegisterFlags for commonCfg
func (c *commonCfg) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(
		&c.defaultAccount,
		"account",
		defaultCommonOptions.defaultAccount,
		"default local account to use",
	)

	fs.StringVar(
		&c.defaultRealm,
		"default-realm",
		defaultCommonOptions.defaultRealm,
		"default realm to display when gnobro starts",
	)

	fs.StringVar(
		&c.sshListener,
		"ssh",
		defaultCommonOptions.sshListener,
		"ssh server listener address",
	)

	fs.StringVar(
		&c.sshHostKeyPath,
		"ssh-key",
		defaultCommonOptions.sshHostKeyPath,
		"ssh host key path",
	)

	fs.BoolVar(
		&c.banner,
		"banner",
		defaultCommonOptions.banner,
		"display a banner",
	)

	fs.BoolVar(
		&c.readonly,
		"readonly",
		defaultCommonOptions.readonly,
		"readonly mode, no commands allowed",
	)

	fs.BoolVar(
		&c.jsonlog,
		"jsonlog",
		defaultCommonOptions.jsonlog,
		"display server log as json format",
	)
}

// execRemote handles the default remote connection mode
func execRemote(ctx context.Context, cfg *remoteCfg, args []string, cio commands.IO) error {
	// Apply defaults if not set
	if cfg.remote == "" {
		cfg.remote = defaultRemoteOptions.remote
	}
	if cfg.chainID == "" {
		cfg.chainID = defaultRemoteOptions.chainID
	}

	// Get node status to display version
	cl, err := client.NewHTTPClient(cfg.remote)
	if err != nil {
		return fmt.Errorf("unable to create http client for %q: %w", cfg.remote, err)
	}

	status, err := cl.Status()
	if err != nil {
		cio.ErrPrintfln("Warning: unable to get node status: %v", err)
	} else {
		cio.Printfln("Connected to %s (version: %s)", cfg.remote, status.NodeInfo.Version)
	}

	// Setup signer and browser
	signer, err := setupSigner(cio, cfg.defaultAccount, cfg.chainID)
	if err != nil {
		return err
	}

	gnocl := &gnoclient.Client{
		RPCClient: cl,
		Signer:    signer,
	}

	// Get realm path from args
	path := getRealmPath(args, cfg.defaultRealm)

	// Create browser config
	bcfg := browser.DefaultConfig()
	bcfg.Readonly = cfg.readonly
	bcfg.URLDefaultValue = path
	bcfg.URLPrefix = gnoPrefix

	return runBrowser(ctx, gnocl, cfg.commonCfg, bcfg, cio, false)
}

// execLocal handles the local gnodev connection mode
func execLocal(ctx context.Context, cfg *localCfg, args []string, cio commands.IO) error {
	// Apply defaults if not set
	if cfg.endpoint == "" {
		cfg.endpoint = defaultLocalOptions.endpoint
	}

	// For now, we still need to know the RPC endpoint
	// TODO: In the future, this should come from the gnodev WebSocket protocol
	rpcEndpoint := "http://127.0.0.1:26657"
	chainID := "dev"

	// Check if gnodev is running
	cl, err := client.NewHTTPClient(rpcEndpoint)
	if err != nil {
		return fmt.Errorf("unable to create http client: %w", err)
	}

	if _, err := cl.Status(); err != nil {
		cio.ErrPrintfln("Error: gnodev doesn't seem to be running")
		cio.ErrPrintfln("")
		cio.ErrPrintfln("To start gnodev, run:")
		cio.ErrPrintfln("  gnodev")
		cio.ErrPrintfln("")
		cio.ErrPrintfln("Or specify a different endpoint with --endpoint")
		return fmt.Errorf("unable to connect to gnodev: %w", err)
	}

	cio.Printfln("Connected to gnodev @ %s", cfg.endpoint)
	cio.Printfln("Hot reload enabled")

	// Setup signer and browser
	signer, err := setupSigner(cio, cfg.defaultAccount, chainID)
	if err != nil {
		return err
	}

	gnocl := &gnoclient.Client{
		RPCClient: cl,
		Signer:    signer,
	}

	// Get realm path from args
	path := getRealmPath(args, cfg.defaultRealm)

	// Create browser config
	bcfg := browser.DefaultConfig()
	bcfg.Readonly = cfg.readonly
	bcfg.URLDefaultValue = path
	bcfg.URLPrefix = gnoPrefix

	// No banner in local mode

	return runBrowserWithDev(ctx, gnocl, cfg.commonCfg, bcfg, cfg.endpoint, cio)
}

// Helper functions

func getRealmPath(args []string, defaultRealm string) string {
	var path string
	if len(args) > 0 {
		path = strings.TrimSpace(args[0])
		path = strings.TrimPrefix(path, gnoPrefix)
	} else if defaultRealm != "" {
		path = strings.TrimLeft(defaultRealm, gnoPrefix)
	}
	return path
}

func setupSigner(io commands.IO, address string, chainID string) (gnoclient.Signer, error) {
	home := gnoenv.HomeDir()
	
	var kb keys.Keybase
	if address != "" {
		var err error
		kb, err = keys.NewKeyBaseFromDir(home)
		if err != nil {
			return nil, fmt.Errorf("unable to load keybase: %w", err)
		}
	} else {
		// create an in-memory keybase
		kb = keys.NewInMemory()
		kb.CreateAccount(integration.DefaultAccount_Name, integration.DefaultAccount_Seed, "", "", 0, 0)
		address = integration.DefaultAccount_Name
	}

	return getSignerForAccount(io, address, kb, chainID)
}

func runBrowser(ctx context.Context, gnocl *gnoclient.Client, cfg commonCfg, bcfg browser.Config, io commands.IO, isDev bool) error {
	if cfg.sshListener != "" {
		return runServer(ctx, gnocl, cfg, bcfg, io)
	}

	if cfg.banner && !isDev {
		bcfg.Banner = NewGnoLandBanner()
	}

	return runTUI(ctx, gnocl, bcfg, io)
}

func runBrowserWithDev(ctx context.Context, gnocl *gnoclient.Client, cfg commonCfg, bcfg browser.Config, wsEndpoint string, io commands.IO) error {
	if cfg.sshListener != "" {
		// SSH server mode doesn't support hot reload yet
		io.ErrPrintfln("Warning: SSH server mode doesn't support hot reload events")
		return runServer(ctx, gnocl, cfg, bcfg, io)
	}

	return runTUIWithDev(ctx, gnocl, bcfg, wsEndpoint, io)
}

func runTUI(ctx context.Context, gnocl *gnoclient.Client, bcfg browser.Config, io commands.IO) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	bcfg.Renderer = lipgloss.DefaultRenderer()
	model := browser.New(bcfg, gnocl)
	
	p := tea.NewProgram(model,
		tea.WithContext(ctx),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	io.Println("Bye!")
	return nil
}

func runTUIWithDev(ctx context.Context, gnocl *gnoclient.Client, bcfg browser.Config, wsEndpoint string, io commands.IO) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	bcfg.Renderer = lipgloss.DefaultRenderer()
	model := browser.New(bcfg, gnocl)
	
	p := tea.NewProgram(model,
		tea.WithContext(ctx),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	var eg errgroup.Group

	// Parse WebSocket endpoint
	devpoint, err := url.Parse(wsEndpoint)
	if err != nil {
		return fmt.Errorf("unable to parse dev endpoint: %w", err)
	}
	devpoint.Path = "_events"

	// Setup WebSocket connection for events
	var devcl browser.DevClient
	devcl.Handler = func(typ events.Type, data any) error {
		switch typ {
		case events.EvtReload, events.EvtReset, events.EvtTxResult:
			p.Send(browser.RefreshRealm())
		default:
			// TODO: Send event to sidebar when implemented
		}
		return nil
	}

	eg.Go(func() error {
		defer cancel()
		if err := devcl.Run(ctx, devpoint.String(), nil); err != nil {
			return fmt.Errorf("dev connection failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		defer cancel()
		_, err := p.Run()
		return err
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	io.Println("Bye!")
	return nil
}


// Keep existing helper functions below...

func runServer(ctx context.Context, gnocl *gnoclient.Client, cfg commonCfg, bcfg browser.Config, io commands.IO) error {
	// setup logger
	logger := newLogger(io.Out(), cfg.jsonlog)

	teaHandler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		shortid := fmt.Sprintf("%.10s", s.Context().SessionID())

		bcfgCopy := bcfg // copy config

		bcfgCopy.Logger = logger.WithGroup(shortid)
		bcfgCopy.Renderer = bubbletea.MakeRenderer(s)

		if cfg.banner {
			bcfgCopy.Banner = NewGnoLandBanner()
		}

		pval := s.Context().Value("path")
		if path, ok := pval.(string); ok && len(path) > 0 {
			// Erase banner on specifc command
			bcfgCopy.Banner = browser.ModelBanner{}
			// Set up url
			bcfgCopy.URLDefaultValue = path
		}

		bcfgCopy.Logger.Info("session started",
			"time", time.Now(),
			"path", bcfgCopy.URLDefaultValue,
			"sid", s.Context().SessionID(),
			"user", s.User())
		model := browser.New(bcfgCopy, gnocl)

		return model, []tea.ProgramOption{
			tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
		}
	}

	sshaddr, err := net.ResolveTCPAddr("", cfg.sshListener)
	if err != nil {
		return fmt.Errorf("unable to resolve address: %w", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(sshaddr.String()),
		wish.WithHostKeyPath(cfg.sshHostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // ensure PTY
			ValidatePathCommandMiddleware(bcfg.URLPrefix),
			StructuredMiddlewareWithLogger(
				ctx, logger, slog.LevelInfo,
			),
			// XXX: add ip throttler
		),
	)

	var errgs errgroup.Group

	errgs.Go(func() error {
		logger.Info("starting SSH server", "addr", sshaddr.String())
		return s.ListenAndServe()
	})

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	errgs.Go(func() error {
		<-ctx.Done()

		logger.Info("stopping SSH server... (5s timeout)")

		sctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		return s.Shutdown(sctx)
	})

	if err := errgs.Wait(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		return err
	}

	if !cfg.jsonlog {
		io.Println("Bye!")
	}
	return nil
}

func getSignerForAccount(io commands.IO, address string, kb keys.Keybase, chainID string) (gnoclient.Signer, error) {
	var signer gnoclient.SignerFromKeybase

	signer.Keybase = kb
	signer.Account = address
	signer.ChainID = chainID

	if ok, err := kb.HasByNameOrAddress(address); !ok || err != nil {
		if err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}

		return nil, fmt.Errorf("unknown name/address: %q", address)
	}

	// try empty password first
	if _, err := kb.ExportPrivKey(address, ""); err != nil {
		prompt := fmt.Sprintf("[%.10s] Enter password:", address)
		signer.Password, err = io.GetPassword(prompt, true)
		if err != nil {
			return nil, fmt.Errorf("error while reading password: %w", err)
		}

		if _, err := kb.ExportPrivKey(address, signer.Password); err != nil {
			return nil, fmt.Errorf("invalid password: %w", err)
		}
	}

	return signer, nil
}

func ValidatePathCommandMiddleware(pathPrefix string) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			switch cmd := s.Command(); len(cmd) {
			case 0: // ok
				next(s)
				return
			case 1: // check for valid path
				path := cmd[0]
				if strings.HasPrefix(path, pathPrefix) && gopath.Clean(path) == path {
					s.Context().SetValue("path", path)
					next(s)
					return
				}

				fmt.Fprintln(s.Stderr(), "provided path is invalid")
			default:
				fmt.Fprintln(s.Stderr(), "too many arguments")
			}

			s.Exit(1)
		}
	}
}

func StructuredMiddlewareWithLogger(ctx context.Context, logger *slog.Logger, level slog.Level) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			ct := time.Now()
			hpk := sess.PublicKey() != nil
			pty, _, _ := sess.Pty()
			logger.Log(
				ctx,
				level,
				"connect",
				"user", sess.User(),
				"remote-addr", sess.RemoteAddr().String(),
				"public-key", hpk,
				"command", sess.Command(),
				"term", pty.Term,
				"width", pty.Window.Width,
				"height", pty.Window.Height,
				"client-version", sess.Context().ClientVersion(),
			)
			next(sess)
			logger.Log(
				ctx,
				level,
				"disconnect",
				"user", sess.User(),
				"remote-addr", sess.RemoteAddr().String(),
				"duration", time.Since(ct),
			)
		}
	}
}

func newLogger(out io.Writer, json bool) *slog.Logger {
	if json {
		return slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	charmlogger := charmlog.New(out)
	charmlogger.SetLevel(charmlog.DebugLevel)
	return slog.New(charmlogger)
}