package gnoland

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	vmm "github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/pkg/gnomod"
	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/bft/types"
	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	osm "github.com/gnolang/gno/tm2/pkg/os"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// LoadGenesisBalancesFile loads genesis balances from the provided file path.
func LoadGenesisBalancesFile(path string) ([]Balance, error) {
	// each balance is in the form: g1xxxxxxxxxxxxxxxx=100000ugnot
	content := osm.MustReadFile(path)
	lines := strings.Split(string(content), "\n")

	balances := make([]Balance, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// remove comments.
		line = strings.Split(line, "#")[0]
		line = strings.TrimSpace(line)

		// skip empty lines.
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=") // <address>=<coin>
		if len(parts) != 2 {
			return nil, errors.New("invalid genesis_balance line: " + line)
		}

		addr, err := crypto.AddressFromBech32(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid balance addr %s: %w", parts[0], err)
		}

		coins, err := std.ParseCoins(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid balance coins %s: %w", parts[1], err)
		}

		balances = append(balances, Balance{
			Address: addr,
			Amount:  coins,
		})
	}

	return balances, nil
}

// LoadGenesisTxsFile loads genesis transactions from the provided file path.
// XXX: Improve the way we generate and load this file
func LoadGenesisTxsFile(path string, chainID string, genesisRemote string) ([]std.Tx, error) {
	txs := []std.Tx{}
	txsBz := osm.MustReadFile(path)
	txsLines := strings.Split(string(txsBz), "\n")
	for _, txLine := range txsLines {
		if txLine == "" {
			continue // Skip empty line.
		}

		// Patch the TX.
		txLine = strings.ReplaceAll(txLine, "%%CHAINID%%", chainID)
		txLine = strings.ReplaceAll(txLine, "%%REMOTE%%", genesisRemote)

		var tx std.Tx
		if err := amino.UnmarshalJSON([]byte(txLine), &tx); err != nil {
			return nil, fmt.Errorf("unable to Unmarshall txs file: %w", err)
		}

		txs = append(txs, tx)
	}

	return txs, nil
}

// LoadPackagesFromDir loads gno packages from a directory.
// It creates and returns a list of transactions based on these packages.
func LoadPackagesFromDir(dir string, creator bft.Address, fee std.Fee, tplData GenesisTplData) ([]std.Tx, error) {
	// list all packages from target path
	pkgs, err := gnomod.ListPkgs(dir) // XXX: gnomod isn't maanaging packages but modules
	if err != nil {
		return nil, fmt.Errorf("listing gno packages: %w", err)
	}

	// Sort packages by dependencies.
	sortedPkgs, err := pkgs.Sort()
	if err != nil {
		return nil, fmt.Errorf("sorting packages: %w", err)
	}

	// Filter out draft packages.
	nonDraftPkgs := sortedPkgs.GetNonDraftPkgs()
	txs := []std.Tx{}
	for _, pkg := range nonDraftPkgs {
		memPkg := gno.ReadMemPackage(pkg.Dir, pkg.Name, tplData)
		tx, err := LoadPackage(memPkg, creator, fee, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to load package %q: %w", pkg.Dir, err)
		}

		txs = append(txs, tx)
	}

	return txs, nil
}

// LoadPackage loads a single package into a `std.Tx`
func LoadPackage(memPkg *std.MemPackage, creator bft.Address, fee std.Fee, deposit std.Coins, tplData GenesisTplData) (std.Tx, error) {
	var tx std.Tx

	// Open files in directory as MemPackage.
	err := memPkg.Validate()
	if err != nil {
		return tx, fmt.Errorf("invalid package: %w", err)
	}

	// Codegen (.tpl)
	for _, file := range memPkg.Files {
		if strings.HasSuffix(file.Name, ".tpl") {
			file.Name = strings.TrimSuffix(file.Name, ".tpl")
			var tpl = template.Must(template.New(file.Name).Parse(file.Body))
			var buf bytes.Buffer
			if err := tpl.Execute(&buf, tplData); err != nil {
				return std.Tx{}, err
			}
			file.Body = buf.String()
			if strings.HasSuffix(file.Name, ".gno") {
				file.Body = "// Code generated by gnoland; DO NOT EDIT.\n" + file.Body
			}
		}
	}

	// Create transaction
	tx.Fee = fee
	tx.Msgs = []std.Msg{
		vmm.MsgAddPackage{
			Creator: creator,
			Package: memPkg,
			Deposit: deposit,
		},
	}
	tx.Signatures = make([]std.Signature, len(tx.GetSigners()))

	return tx, nil
}

type GenesisTplData struct {
	Genesis types.GenesisDoc
	// LocalGnokey
	// BuildTags
	// -X
}
