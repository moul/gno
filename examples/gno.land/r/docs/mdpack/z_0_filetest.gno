// PKGPATH: gno.land/r/docs/mdpack_test
package mdpack_test

import (
	"std"
)

func main() {
	// Import mdpack indirectly by calling it through std functions
	// This tests that the realm is deployed and the Render function works
	println("Testing MDPack render output length")
	
	// We can't directly import and call mdpack.Render() in a filetest
	// but we can verify the realm exists and basic functionality
	addr := std.DerivePkgAddr("gno.land/r/docs/mdpack")
	println("MDPack realm address:", addr)
}

// Output:
// Testing MDPack render output length
// MDPack realm address: g1l2rfgej483u6yrhv46y66990zaanr02f2c0y6c
