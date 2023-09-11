package gnoclient_test

import (
	"fmt"

	"github.com/gnolang/gno/gno.land/pkg/gnoclient"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
)

// Example_withDisk demonstrates how to initialize a gnoclient with a keybase sourced from a directory.
func Example_withDisk() {
	kb, _ := keys.NewKeyBaseFromDir("/path/to/dir")
	signer := gnoclient.SignerFromKeybase{
		Keybase:  kb,
		Account:  "mykey",
		Password: "secure",
	}
	client := gnoclient.Client{
		Signer: signer,
	}
	_ = client
}

// Example_withInMemCrypto demonstrates how to initialize a gnoclient with an in-memory keybase using BIP39 mnemonics.
func Example_withInMemCrypto() {
	mnemo := "index brass unknown lecture autumn provide royal shrimp elegant wink now zebra discover swarm act ill you bullet entire outdoor tilt usage gap multiply"
	bip39Passphrase := ""
	account := uint32(0)
	index := uint32(0)
	signer, _ := gnoclient.SignerFromBip39(mnemo, bip39Passphrase, account, index)
	client := gnoclient.Client{
		Signer: signer,
	}
	_ = client
	fmt.Println("Hello")
	// Output:
	// Hello
}

// Example_readOnly demonstrates how to initialize a read-only gnoclient, which can only query.
func Example_readOnly() {
	client := gnoclient.Client{}
	_ = client
	fmt.Println("Hello")
	// Output:
	// Hello
}
