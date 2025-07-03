package dao_sudo

import "std"

// Message types that implement the realm.Message interface
// These would normally be imported from the SDK, but we define them here
// for the example.

// MsgAddPackage represents a message to deploy a new package/realm
type MsgAddPackage struct {
	Creator std.Address
	Package Package
	Deposit std.Coins
}

func (m *MsgAddPackage) Route() string { return "vm" }
func (m *MsgAddPackage) Type() string  { return "add_package" }
func (m *MsgAddPackage) ValidateBasic() error {
	if m.Creator.IsZero() {
		return std.ErrInvalidAddress("empty creator address")
	}
	if m.Package.Path == "" {
		return std.ErrInvalidRequest("empty package path")
	}
	return nil
}
func (m *MsgAddPackage) GetSignBytes() []byte {
	// In real implementation, this would serialize the message
	return []byte(m.Package.Path)
}

// MsgCall represents a message to call a realm function
type MsgCall struct {
	Caller  std.Address
	PkgPath string
	Func    string
	Args    []string
	Send    std.Coins
}

func (m *MsgCall) Route() string { return "vm" }
func (m *MsgCall) Type() string  { return "call" }
func (m *MsgCall) ValidateBasic() error {
	if m.Caller.IsZero() {
		return std.ErrInvalidAddress("empty caller address")
	}
	if m.PkgPath == "" {
		return std.ErrInvalidRequest("empty package path")
	}
	if m.Func == "" {
		return std.ErrInvalidRequest("empty function name")
	}
	return nil
}
func (m *MsgCall) GetSignBytes() []byte {
	// In real implementation, this would serialize the message
	return []byte(m.PkgPath + "." + m.Func)
}

// Package represents a Gno package
type Package struct {
	Name  string
	Path  string
	Files []MemFile
}

// MemFile represents a file in memory
type MemFile struct {
	Name string
	Body string
}

// MsgSend represents a bank send message
type MsgSend struct {
	FromAddress std.Address
	ToAddress   std.Address
	Amount      std.Coins
}

func (m *MsgSend) Route() string { return "bank" }
func (m *MsgSend) Type() string  { return "send" }
func (m *MsgSend) ValidateBasic() error {
	if m.FromAddress.IsZero() {
		return std.ErrInvalidAddress("empty from address")
	}
	if m.ToAddress.IsZero() {
		return std.ErrInvalidAddress("empty to address")
	}
	if !m.Amount.IsValid() {
		return std.ErrInvalidCoins("invalid amount")
	}
	return nil
}
func (m *MsgSend) GetSignBytes() []byte {
	return []byte(m.FromAddress.String() + m.ToAddress.String())
}

// MsgDelegate represents a staking delegation message
type MsgDelegate struct {
	DelegatorAddress std.Address
	ValidatorAddress std.Address
	Amount           std.Coins
}

func (m *MsgDelegate) Route() string { return "staking" }
func (m *MsgDelegate) Type() string  { return "delegate" }
func (m *MsgDelegate) ValidateBasic() error {
	if m.DelegatorAddress.IsZero() {
		return std.ErrInvalidAddress("empty delegator address")
	}
	if m.ValidatorAddress.IsZero() {
		return std.ErrInvalidAddress("empty validator address")
	}
	if !m.Amount.IsValid() {
		return std.ErrInvalidCoins("invalid amount")
	}
	return nil
}
func (m *MsgDelegate) GetSignBytes() []byte {
	return []byte(m.DelegatorAddress.String() + m.ValidatorAddress.String())
}