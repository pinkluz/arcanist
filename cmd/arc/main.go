package main

import (
	"github.com/pinkluz/arcanist/cli"

	_ "github.com/pinkluz/arcanist/cli/cascade"
	_ "github.com/pinkluz/arcanist/cli/flow"
)

func main() {
	cli.GetRoot().Execute()
}
