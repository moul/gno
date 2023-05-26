// PKGPATH: gno.land/r/demo/tests_test
package tests_test

var amap map[string]string = map[string]string{"a": "1", "b": "2", "c": "3"}

func main() {
	println("premain", amap)
	amap["d"] = "4"
	amap["b"] = "22"
	amap["e"] = "5"
	println("postmain", amap)
}

// Output:
// premain map{}
// postmain map{("d" string):("4" string),("b" string):("22" string),("e" string):("5" string)}

