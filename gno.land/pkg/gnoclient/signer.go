package gnoclient

import "github.com/gnolang/gno/tm2/pkg/crypto/keys"

// Signer provides an interface for signing.
type Signer interface {
	Sign()           // TBD
	Info() keys.Info // returns keys info, containing the address.
	Validate() error // checks wether the signer is well configured.
}

// SignerFromKeybase represents a signer created from a Keybase.
type SignerFromKeybase struct {
	Keybase  keys.Keybase // Holds keys in memory or on disk
	Account  string       // Could be name or bech32 format
	Password string       // Password for encryption
}

func (s SignerFromKeybase) Validate() error {
	_, err := s.Keybase.GetByNameOrAddress(s.Account)
	if err != nil {
		return err
	}
	// XXX: additional checks?
	return nil
}

func (s SignerFromKeybase) Info() keys.Info {
	info, err := s.Keybase.GetByNameOrAddress(s.Account)
	if err != nil {
		panic("should not happen")
	}
	return info
}

// Sign implements the Signer interface for SignerFromKeybase.
func (s SignerFromKeybase) Sign() {
	panic("not implemented")
}

// Ensure SignerFromKeybase implements Signer interface.
var _ Signer = (*SignerFromKeybase)(nil)

// SignerFromBip39 creates an in-memory keybase with a single default account.
// This can be useful in scenarios where storing private keys in the filesystem isn't feasible.
//
// Warning: Using keys.NewKeyBaseFromDir is recommended where possible, as it is more secure.
func SignerFromBip39(mnemo string, passphrase string, account uint32, index uint32) (Signer, error) {
	kb := keys.NewInMemory()
	name := "default"
	passwd := "" // Password isn't needed for in-memory storage

	_, err := kb.CreateAccount(name, mnemo, passphrase, passwd, account, index)
	if err != nil {
		return nil, err
	}

	signer := SignerFromKeybase{
		Keybase:  kb,
		Account:  name,
		Password: passwd,
	}

	return &signer, nil
}
