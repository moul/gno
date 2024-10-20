package vm

import (
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/sdk/params"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// ----------------------------------------
// SDKBanker

type SDKBanker struct {
	vmk *VMKeeper
	ctx sdk.Context
}

func NewSDKBanker(vmk *VMKeeper, ctx sdk.Context) *SDKBanker {
	return &SDKBanker{
		vmk: vmk,
		ctx: ctx,
	}
}

func (bnk *SDKBanker) GetCoins(b32addr crypto.Bech32Address) (dst std.Coins) {
	addr := crypto.MustAddressFromString(string(b32addr))
	coins := bnk.vmk.bank.GetCoins(bnk.ctx, addr)
	return coins
}

func (bnk *SDKBanker) SendCoins(b32from, b32to crypto.Bech32Address, amt std.Coins) {
	from := crypto.MustAddressFromString(string(b32from))
	to := crypto.MustAddressFromString(string(b32to))
	err := bnk.vmk.bank.SendCoins(bnk.ctx, from, to, amt)
	if err != nil {
		panic(err)
	}
}

func (bnk *SDKBanker) TotalCoin(denom string) int64 {
	panic("not yet implemented")
}

func (bnk *SDKBanker) IssueCoin(b32addr crypto.Bech32Address, denom string, amount int64) {
	addr := crypto.MustAddressFromString(string(b32addr))
	_, err := bnk.vmk.bank.AddCoins(bnk.ctx, addr, std.Coins{std.Coin{Denom: denom, Amount: amount}})
	if err != nil {
		panic(err)
	}
}

func (bnk *SDKBanker) RemoveCoin(b32addr crypto.Bech32Address, denom string, amount int64) {
	addr := crypto.MustAddressFromString(string(b32addr))
	_, err := bnk.vmk.bank.SubtractCoins(bnk.ctx, addr, std.Coins{std.Coin{Denom: denom, Amount: amount}})
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------
// SDKParams

type SDKParams struct {
	vmk *VMKeeper
	ctx sdk.Context
}

func NewSDKParams(vmk *VMKeeper, ctx sdk.Context) *SDKParams {
	return &SDKParams{
		vmk: vmk,
		ctx: ctx,
	}
}

// SetXXX helpers:
// - dynamically register a new key with the corresponding type in the paramset table (only once).
// - set the value.

func (prm *SDKParams) SetString(key, value string) {
	// if !prm.vmk.prmk.Has(prm.ctx, key) {
	// XXX: bad workaround, maybe we should have a dedicated "dynamic keeper" allowing to create keys on the go?
	if !prm.vmk.prmk.HasTypeKey(key) {
		prm.vmk.prmk.RegisterType(params.NewParamSetPair(key, "", validateNoOp))
	}
	prm.vmk.prmk.Set(prm.ctx, key, value)
}

func (prm *SDKParams) SetBool(key string, value bool) {
	// if !prm.vmk.prmk.Has(prm.ctx, key) {
	if !prm.vmk.prmk.HasTypeKey(key) {
		prm.vmk.prmk.RegisterType(params.NewParamSetPair(key, true, validateNoOp))
	}
	prm.vmk.prmk.Set(prm.ctx, key, value)
}

func (prm *SDKParams) SetInt64(key string, value int64) {
	// if !prm.vmk.prmk.Has(prm.ctx, key) {
	if !prm.vmk.prmk.HasTypeKey(key) {
		prm.vmk.prmk.RegisterType(params.NewParamSetPair(key, int64(0), validateNoOp))
	}
	prm.vmk.prmk.Set(prm.ctx, key, value)
}

func (prm *SDKParams) SetUint64(key string, value uint64) {
	// if !prm.vmk.prmk.Has(prm.ctx, key) {
	if !prm.vmk.prmk.HasTypeKey(key) {
		prm.vmk.prmk.RegisterType(params.NewParamSetPair(key, uint64(0), validateNoOp))
	}
	prm.vmk.prmk.Set(prm.ctx, key, value)
}

func validateNoOp(_ interface{}) error { return nil }
