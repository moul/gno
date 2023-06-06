// PKGPATH: gno.land/r/demo/tests_test
package tests_test

import "os"
import "fmt"

var amap map[string]string = map[string]string{"a": "1", "b": "2", "c": "3"}

func main() {
	println("premain", amap)
	amap["d"] = "4"
	amap["b"] = "22"
	amap["f"] = "22"
	amap["g"] = "22"
	amap["h"] = "22"
	amap["i"] = "22"
	amap["j"] = "22"
	amap["g"] = "22"
	amap["b"] = "22"
	amap["b"] = "22"
	amap["e"] = "5"
	println("postmain", amap)
}

// Output:
// premain map{}
// postmain map{("d" string):("4" string),("b" string):("22" string),("f" string):("22" string),("g" string):("22" string),("h" string):("22" string),("i" string):("22" string),("j" string):("22" string),("e" string):("5" string)}

