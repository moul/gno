package valopers

import (
	"std"

	"gno.land/p/moul/authz"
)

var (
	auth       authz.Authority
	privConfig *authz.PBACPrivateConfig
)

func init() {
	// Initialize with origin caller as admin
	auth, privConfig = authz.NewWithOrigin()
	
	// Origin is automatically set as admin
	// Admin has all permissions by default
}

func Auth() authz.Authority {
	return auth
}

func updateInstructions(newInstructions string) {
	err := auth.DoByCurrent("update-instructions", func() error {
		instructions = newInstructions
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func updateMinFee(newMinFee int64) {
	err := auth.DoByCurrent("update-min-fee", func() error {
		minFee = std.NewCoin("ugnot", newMinFee)
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func NewInstructionsProposalCallback(newInstructions string) func() error {
	cb := func() error {
		updateInstructions(newInstructions)
		return nil
	}

	return cb
}

func NewMinFeeProposalCallback(newMinFee int64) func() error {
	cb := func() error {
		updateMinFee(newMinFee)
		return nil
	}

	return cb
}
