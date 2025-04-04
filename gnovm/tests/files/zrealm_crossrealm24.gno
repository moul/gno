// PKGPATH: gno.land/r/crossrealm_test
package crossrealm_test

var counter = 42

func main() {
	switchrealm()
	//switchrealm()
	// println(crossrealm.Notebook{Counter: &counter}.Inc().String())
	// println(crossrealm.Notebook{Counter: &crossrealm.Counter}.Inc().String()) // cannot modify external-realm or non-realm object
	//print(".")
}

// Output:
//
