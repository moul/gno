package realm

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs"
)

func init() {
	stdlibs.Register("realm", "sudo", X_sudo)
}

var Package = &gno.PackageNode{
	Name: "realm",
	Path: "realm",
	Files: []*gno.FileNode{
		{
			Name: "realm.gno",
			Body: `
package realm

type Message interface {
	Route() string
	Type() string
	ValidateBasic() error
	GetSignBytes() []byte
}

func Sudo(msg Message) {
	sudo(msg)
}

func sudo(msg Message) // native
`,
		},
	},
}