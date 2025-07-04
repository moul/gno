package grc20reg

import (
	"std"

	"gno.land/p/demo/avl"
	"gno.land/p/demo/avl/rotree"
	"gno.land/p/demo/fqname"
	"gno.land/p/demo/grc/grc20"
	"gno.land/p/demo/ufmt"
	"gno.land/p/moul/mdtable"
)

var registry = avl.NewTree() // rlmPath[.slug] -> *Token (slug is optional)
func Register(cur realm, token *grc20.Token, slug string) {
	rlmPath := std.PreviousRealm().PkgPath()
	key := fqname.Construct(rlmPath, slug)
	registry.Set(key, token)
	std.Emit(
		registerEvent,
		"pkgpath", rlmPath,
		"slug", slug,
	)
}

func Get(key string) *grc20.Token {
	token, ok := registry.Get(key)
	if !ok {
		return nil
	}
	return token.(*grc20.Token)
}

func MustGet(key string) *grc20.Token {
	token := Get(key)
	if token == nil {
		panic("unknown token: " + key)
	}
	return token
}

func Render(path string) string {
	switch {
	case path == "": // home
		// TODO: add pagination
		s := "# GRC20 Token Registry\n\n"
		
		// Create table for token listings
		table := mdtable.New([]string{"Token Name", "Symbol", "Decimals", "Realm", "Info"})
		count := 0
		
		registry.Iterate("", "", func(key string, tokenI any) bool {
			count++
			token := tokenI.(*grc20.Token)
			rlmPath, slug := fqname.Parse(key)
			rlmLink := fqname.RenderLink(rlmPath, slug)
			infoLink := ufmt.Sprintf("[details](/r/demo/grc20reg:%s)", key)
			
			table.AddRow([]string{
				token.GetName(),
				token.GetSymbol(),
				ufmt.Sprintf("%d", token.GetDecimals()),
				rlmLink,
				infoLink,
			})
			return false
		})
		
		if count == 0 {
			return "No registered tokens."
		}
		
		s += table.String()
		return s
	default: // specific token
		key := path
		token := MustGet(key)
		rlmPath, slug := fqname.Parse(key)
		rlmLink := fqname.RenderLink(rlmPath, slug)
		s := ufmt.Sprintf("# %s\n", token.GetName())
		s += ufmt.Sprintf("- symbol: **%s**\n", token.GetSymbol())
		s += ufmt.Sprintf("- realm: %s\n", rlmLink)
		s += ufmt.Sprintf("- decimals: %d\n", token.GetDecimals())
		s += ufmt.Sprintf("- total supply: %d\n", token.TotalSupply())
		return s
	}
}

const registerEvent = "register"

func GetRegistry() *rotree.ReadOnlyTree {
	return rotree.Wrap(registry, nil)
}
