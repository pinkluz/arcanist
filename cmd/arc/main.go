package main

import (
	"github.com/pinkluz/arcanist/cli"

	_ "github.com/pinkluz/arcanist/cli/cascade"
	_ "github.com/pinkluz/arcanist/cli/flow"
	_ "github.com/pinkluz/arcanist/cli/graft"
	_ "github.com/pinkluz/arcanist/cli/prune"
)

func main() {
	cli.GetRoot().Execute()
}
