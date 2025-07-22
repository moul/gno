package vm

import (
	"time"

	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// Route returns the route string for session messages
const RouterKeySession = "vm"

//----------------------------------------
// MsgCreateSession

// MsgCreateSession creates a new session for an account
type MsgCreateSession struct {
	Creator           crypto.Address `json:"creator" yaml:"creator"`
	SessionKey        crypto.PubKey  `json:"session_key" yaml:"session_key"`
	TransferCapacity  std.Coins      `json:"transfer_capacity" yaml:"transfer_capacity"`
	UnlimitedTransfer bool           `json:"unlimited_transfer" yaml:"unlimited_transfer"`
	CanManageSessions bool           `json:"can_manage_sessions" yaml:"can_manage_sessions"`
	CanManagePackages bool           `json:"can_manage_packages" yaml:"can_manage_packages"`
	ValidationOnly    bool           `json:"validation_only" yaml:"validation_only"`
	ExpirationTime    time.Time      `json:"expiration_time" yaml:"expiration_time"`
	RealmsWhitelist   []string       `json:"realms_whitelist" yaml:"realms_whitelist"`
	InitialSequence   uint64         `json:"initial_sequence" yaml:"initial_sequence"`
}

var _ sdk.Msg = MsgCreateSession{}

// NewMsgCreateSession creates a new MsgCreateSession instance
func NewMsgCreateSession(
	creator crypto.Address,
	sessionKey crypto.PubKey,
) MsgCreateSession {
	return MsgCreateSession{
		Creator:    creator,
		SessionKey: sessionKey,
	}
}

// Implements sdk.Msg
func (msg MsgCreateSession) Route() string { return RouterKeySession }

// Implements sdk.Msg
func (msg MsgCreateSession) Type() string { return "create_session" }

// Implements sdk.Msg
func (msg MsgCreateSession) ValidateBasic() error {
	if msg.Creator.IsZero() {
		return std.ErrInvalidAddress("missing creator address")
	}
	if msg.SessionKey == nil {
		return std.ErrInvalidPubKey("missing session public key")
	}
	if !msg.TransferCapacity.IsValid() {
		return std.ErrInvalidCoins("invalid transfer capacity")
	}
	return nil
}

// Implements sdk.Msg
func (msg MsgCreateSession) GetSignBytes() []byte {
	bz, err := amino.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// Implements sdk.Msg
func (msg MsgCreateSession) GetSigners() []crypto.Address {
	return []crypto.Address{msg.Creator}
}

//----------------------------------------
// MsgRevokeSession

// MsgRevokeSession revokes a session from an account
type MsgRevokeSession struct {
	Creator    crypto.Address `json:"creator" yaml:"creator"`
	SessionKey crypto.PubKey  `json:"session_key" yaml:"session_key"`
}

var _ sdk.Msg = MsgRevokeSession{}

// NewMsgRevokeSession creates a new MsgRevokeSession instance
func NewMsgRevokeSession(
	creator crypto.Address,
	sessionKey crypto.PubKey,
) MsgRevokeSession {
	return MsgRevokeSession{
		Creator:    creator,
		SessionKey: sessionKey,
	}
}

// Implements sdk.Msg
func (msg MsgRevokeSession) Route() string { return RouterKeySession }

// Implements sdk.Msg
func (msg MsgRevokeSession) Type() string { return "revoke_session" }

// Implements sdk.Msg
func (msg MsgRevokeSession) ValidateBasic() error {
	if msg.Creator.IsZero() {
		return std.ErrInvalidAddress("missing creator address")
	}
	if msg.SessionKey == nil {
		return std.ErrInvalidPubKey("missing session public key")
	}
	return nil
}

// Implements sdk.Msg
func (msg MsgRevokeSession) GetSignBytes() []byte {
	bz, err := amino.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// Implements sdk.Msg
func (msg MsgRevokeSession) GetSigners() []crypto.Address {
	return []crypto.Address{msg.Creator}
}