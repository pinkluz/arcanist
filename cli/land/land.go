package land

import (
	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
)

type landCmd struct {
	cmd *cobra.Command
}

func (f *landCmd) run(cmd *cobra.Command, args []string) {
	// TODO: do this
}

func init() {
	land := &landCmd{}
	land.cmd = &cobra.Command{
		Use:   "land",
		Short: "Move all branches from the current branch to the parent",
		Long: `When land is run it takes all commits on the current branch and squashes them down to one commit
	and applys it to the parent. It will then try to push the branch if it's remote is not local.`,
		Run: land.run,
	}

	cli.GetRoot().AddCommand(land.cmd)
}
