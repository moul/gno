package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Config tests ---

func TestConfig_SetGet(t *testing.T) {
	home := t.TempDir()

	cfg := &Config{}
	require.NoError(t, ConfigSet(cfg, "key", "moul"))
	require.NoError(t, ConfigSet(cfg, "gas-buffer", "30"))
	require.NoError(t, SaveConfig(home, cfg))

	loaded, err := LoadConfig(home)
	require.NoError(t, err)
	assert.Equal(t, "moul", loaded.Key)
	assert.Equal(t, 30, loaded.GasBuffer)

	val, err := ConfigGet(loaded, "key")
	require.NoError(t, err)
	assert.Equal(t, "moul", val)

	val, err = ConfigGet(loaded, "gas-buffer")
	require.NoError(t, err)
	assert.Equal(t, "30", val)
}

func TestConfig_Defaults(t *testing.T) {
	home := t.TempDir()

	// No config file → defaults
	cfg, err := LoadConfig(home)
	require.NoError(t, err)
	assert.Equal(t, "", cfg.Key)
	assert.Equal(t, defaultGasBuffer, cfg.GetGasBuffer())
}

func TestConfig_UnknownKey(t *testing.T) {
	cfg := &Config{}
	_, err := ConfigGet(cfg, "nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown config key")

	err = ConfigSet(cfg, "nonexistent", "value")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown config key")
}

func TestConfig_GasBufferValidation(t *testing.T) {
	cfg := &Config{}
	err := ConfigSet(cfg, "gas-buffer", "notanumber")
	require.Error(t, err)

	err = ConfigSet(cfg, "gas-buffer", "-1")
	require.Error(t, err)

	require.NoError(t, ConfigSet(cfg, "gas-buffer", "0"))
	// 0 → falls back to default
	assert.Equal(t, defaultGasBuffer, cfg.GetGasBuffer())
}

func TestConfig_List(t *testing.T) {
	cfg := &Config{Key: "alice", GasBuffer: 25}
	list := ConfigList(cfg)
	assert.Contains(t, list, "key=alice")
	assert.Contains(t, list, "gas-buffer=25")
}

func TestConfigCmd_SetGetList(t *testing.T) {
	home := t.TempDir()
	base := &baseCfg{home: home}

	makeIO := func() (commands.IO, *bytes.Buffer) {
		var buf bytes.Buffer
		io := commands.NewTestIO()
		io.SetOut(commands.WriteNopCloser(&buf))
		return io, &buf
	}

	// config set key=bob
	io, buf := makeIO()
	cmd := newConfigCmd(base, io)
	require.NoError(t, cmd.ParseAndRun(context.Background(), []string{"set", "key=bob"}))
	assert.Contains(t, buf.String(), "key=bob")

	// config get key
	io, buf = makeIO()
	cmd = newConfigCmd(base, io)
	require.NoError(t, cmd.ParseAndRun(context.Background(), []string{"get", "key"}))
	assert.Contains(t, buf.String(), "bob")

	// config list
	io, buf = makeIO()
	cmd = newConfigCmd(base, io)
	require.NoError(t, cmd.ParseAndRun(context.Background(), []string{"list"}))
	assert.Contains(t, buf.String(), "key=bob")
}

// --- Version tests ---

func TestVersionCmd(t *testing.T) {
	var outBuf bytes.Buffer
	io := commands.NewTestIO()
	io.SetOut(commands.WriteNopCloser(&outBuf))

	cmd := newVersionCmd(io)
	err := cmd.ParseAndRun(context.Background(), []string{})
	require.NoError(t, err)
	assert.Contains(t, outBuf.String(), "gnopie")
	assert.Contains(t, outBuf.String(), version)
}

// --- Completion tests ---

func TestCompletionCmd(t *testing.T) {
	shells := []struct {
		shell   string
		contain string
	}{
		{"bash", "_gnopie"},
		{"zsh", "#compdef gnopie"},
		{"fish", "complete -c gnopie"},
	}

	for _, tc := range shells {
		t.Run(tc.shell, func(t *testing.T) {
			var outBuf bytes.Buffer
			io := commands.NewTestIO()
			io.SetOut(commands.WriteNopCloser(&outBuf))

			cmd := newCompletionCmd(io)
			err := cmd.ParseAndRun(context.Background(), []string{tc.shell})
			require.NoError(t, err)
			assert.Contains(t, outBuf.String(), tc.contain)
		})
	}
}

func TestCompletionCmd_UnknownShell(t *testing.T) {
	io := commands.NewTestIO()
	cmd := newCompletionCmd(io)
	err := cmd.ParseAndRun(context.Background(), []string{"powershell"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell")
}

func TestCompletionCmd_NoArgs(t *testing.T) {
	io := commands.NewTestIO()
	cmd := newCompletionCmd(io)
	err := cmd.ParseAndRun(context.Background(), []string{})
	require.Error(t, err)
}

// --- ParsePath tests ---

func TestParsePath(t *testing.T) {
	tt := []struct {
		input      string
		wantKind   PathKind
		wantDomain string
		wantPkg    string
		wantSym    string
		wantArgs   []string
		wantFile   string
		wantErr    bool
	}{
		{
			input:   "",
			wantErr: true,
		},
		{
			input:      "gno.land",
			wantKind:   PathNetwork,
			wantDomain: "gno.land",
		},
		{
			input:      "gno.land/r/demo/counter",
			wantKind:   PathPackage,
			wantDomain: "gno.land",
			wantPkg:    "gno.land/r/demo/counter",
		},
		{
			input:      "gno.land/r/demo/counter.Increment",
			wantKind:   PathSymbol,
			wantDomain: "gno.land",
			wantPkg:    "gno.land/r/demo/counter",
			wantSym:    "Increment",
		},
		{
			input:    "gno.land/r/demo/counter.Increment()",
			wantKind: PathCall,
			wantPkg:  "gno.land/r/demo/counter",
			wantSym:  "Increment",
		},
		{
			input:    `gno.land/r/demo/counter.Foo("hello", "world")`,
			wantKind: PathCall,
			wantPkg:  "gno.land/r/demo/counter",
			wantSym:  "Foo",
			wantArgs: []string{"hello", "world"},
		},
		{
			input:    "gno.land/r/demo/counter/counter.gno",
			wantKind: PathFile,
			wantPkg:  "gno.land/r/demo/counter",
			wantFile: "counter.gno",
		},
		{
			input:    "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5",
			wantKind: PathAddress,
		},
		{
			input:      "https://gno.land/r/demo/counter",
			wantKind:   PathPackage,
			wantDomain: "gno.land",
			wantPkg:    "gno.land/r/demo/counter",
		},
		{
			input:      "https://gno.land/r/demo/counter#some-anchor",
			wantKind:   PathPackage,
			wantDomain: "gno.land",
			wantPkg:    "gno.land/r/demo/counter",
		},
		{
			input:      "gno.land/r/demo/counter:some/path",
			wantKind:   PathPackage,
			wantDomain: "gno.land",
			wantPkg:    "gno.land/r/demo/counter",
		},
		{
			input:    "gno.land/u/moul",
			wantKind: PathUser,
			wantSym:  "moul",
		},
		{
			input:    "gno.land/r",
			wantKind: PathNamespace,
			wantPkg:  "gno.land/r",
		},
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			p, err := ParsePath(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantKind, p.Kind)
			if tc.wantDomain != "" {
				assert.Equal(t, tc.wantDomain, p.Domain)
			}
			if tc.wantPkg != "" {
				assert.Equal(t, tc.wantPkg, p.PkgPath)
			}
			if tc.wantSym != "" {
				assert.Equal(t, tc.wantSym, p.Symbol)
			}
			if tc.wantArgs != nil {
				assert.Equal(t, tc.wantArgs, p.Args)
			}
			if tc.wantFile != "" {
				assert.Equal(t, tc.wantFile, p.File)
			}
		})
	}
}

// --- extractDecl tests ---

func TestExtractDecl(t *testing.T) {
	src := `package counter

var counter int

func Increment(cross realm) {
	counter++
}

func Render(_ string) string {
	return ufmt.Sprintf("%d", counter)
}
`
	t.Run("func", func(t *testing.T) {
		decl := extractDecl(src, "Increment")
		assert.Contains(t, decl, "func Increment")
		assert.Contains(t, decl, "counter++")
	})

	t.Run("var", func(t *testing.T) {
		decl := extractDecl(src, "counter")
		assert.Contains(t, decl, "counter int")
	})

	t.Run("not found", func(t *testing.T) {
		decl := extractDecl(src, "DoesNotExist")
		assert.Equal(t, "", decl)
	})
}

// --- extractVarsConsts tests ---

func TestExtractVarsConsts(t *testing.T) {
	src := `package foo

var x int
var y string = "hello"

const MaxItems = 100
const prefix = "gno"

var (
	a int
	b string
)

const (
	Alpha = 1
	Beta  = 2
)
`
	vars, consts := extractVarsConsts(src)
	assert.Contains(t, vars, "x int")
	assert.Contains(t, vars, `y string = "hello"`)
	assert.Contains(t, vars, "a int")
	assert.Contains(t, vars, "b string")
	assert.Contains(t, consts, "MaxItems = 100")
	assert.Contains(t, consts, `prefix = "gno"`)
	assert.Contains(t, consts, "Alpha = 1")
	assert.Contains(t, consts, "Beta  = 2")
}

// --- defaultHome tests ---

func TestDefaultHome_GNOHOME(t *testing.T) {
	t.Setenv("GNOHOME", "/tmp/gnohome_test")
	t.Setenv("GNO_HOME", "")
	assert.Equal(t, "/tmp/gnohome_test", defaultHome())
}

func TestDefaultHome_GnoHome(t *testing.T) {
	t.Setenv("GNOHOME", "")
	t.Setenv("GNO_HOME", "/tmp/gno_home_test")
	assert.Equal(t, "/tmp/gno_home_test", defaultHome())
}

// --- queryCacheKey tests ---

func TestQueryCacheKey(t *testing.T) {
	k1 := queryCacheKey("vm/qfile", "gno.land/r/demo/counter")
	k2 := queryCacheKey("vm/qfile", "gno.land/r/demo/counter")
	k3 := queryCacheKey("vm/qfuncs", "gno.land/r/demo/counter")

	assert.Equal(t, k1, k2, "same inputs → same key")
	assert.NotEqual(t, k1, k3, "different queryPath → different key")
}

func TestLoadCachedQuery_Miss(t *testing.T) {
	home := t.TempDir()
	_, ok := loadCachedQuery(home, "vm/qfile", "nonexistent")
	assert.False(t, ok)
}

func TestSaveLoadCachedQuery(t *testing.T) {
	home := t.TempDir()
	saveCachedQuery(home, "vm/qfile", "gno.land/r/demo/counter", "counter.gno\n")
	result, ok := loadCachedQuery(home, "vm/qfile", "gno.land/r/demo/counter")
	require.True(t, ok)
	assert.Equal(t, "counter.gno\n", result)
}

func TestLoadCachedQuery_Expired(t *testing.T) {
	home := t.TempDir()
	dir := queryCacheDir(home)
	require.NoError(t, os.MkdirAll(dir, 0o755))

	key := queryCacheKey("vm/qfile", "expiredtest")
	path := filepath.Join(dir, key)
	require.NoError(t, os.WriteFile(path, []byte("old data"), 0o644))

	// Backdate the file modification time to 2 hours ago (beyond 1h TTL)
	twoHoursAgo := time.Now().Add(-2 * time.Hour)
	require.NoError(t, os.Chtimes(path, twoHoursAgo, twoHoursAgo))

	_, ok := loadCachedQuery(home, "vm/qfile", "expiredtest")
	assert.False(t, ok, "expired cache entry should not be returned")
}

// --- isNumeric tests ---

func TestIsNumeric(t *testing.T) {
	assert.True(t, isNumeric("0"))
	assert.True(t, isNumeric("123"))
	assert.True(t, isNumeric("-5"))
	assert.False(t, isNumeric(""))
	assert.False(t, isNumeric("abc"))
	assert.False(t, isNumeric("1.5"))
}

// --- joinArgs tests ---

func TestJoinArgs(t *testing.T) {
	assert.Equal(t, `"hello","world"`, joinArgs([]string{"hello", "world"}))
	assert.Equal(t, `42,true`, joinArgs([]string{"42", "true"}))
	assert.Equal(t, `"foo"`, joinArgs([]string{"foo"}))
	assert.Equal(t, ``, joinArgs([]string{}))
}

// --- generateRunCode tests ---

func TestGenerateRunCode(t *testing.T) {
	code := generateRunCode("gno.land/r/demo/counter", "counter", "Increment", []string{"__cross__"})
	assert.Contains(t, code, "package main")
	assert.Contains(t, code, `"gno.land/r/demo/counter"`)
	assert.Contains(t, code, "counter.Increment(cross)")
}

func TestGenerateRunCode_WithArgs(t *testing.T) {
	code := generateRunCode("gno.land/r/demo/foo", "foo", "Bar", []string{"__cross__", "hello"})
	assert.Contains(t, code, `foo.Bar(cross, "hello")`)
}

func TestGenerateRunCode_NoArgs(t *testing.T) {
	code := generateRunCode("gno.land/r/demo/foo", "foo", "Baz", nil)
	assert.Contains(t, code, "foo.Baz()")
}
