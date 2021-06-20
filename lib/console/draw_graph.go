package console

import (
	"fmt"

	"github.com/pinkluz/arcanist/lib/git"
)

var (
	vertical_left  = "\x25\x27"
	vertical_right = "\x25\x1C"
	vertical       = "\x25\x02"
	horizontal     = "\x25\x00"
	up_right       = "\x25\x14"
)

// DrawGraphOpts is left for later when nwe want to allow the user some control
// over the output. For now it's just empty.
type DrawGraphOpts struct {
}

// DrawGraph takes a git.BranchNodeWrapper and renders the output for your
// console. This is all returned as a string so you can
func DrawGraph(bnw git.BranchNodeWrapper, opts *DrawGraphOpts) string {
	internalOpts := opts
	if opts == nil {
		internalOpts = &DrawGraphOpts{}
	}

	fmt.Println(internalOpts)

	fmt.Println(vertical_left, vertical_right, vertical, horizontal, up_right)

	return ""
}
