package prune

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
	"github.com/pinkluz/arcanist/lib/git"
)

type pruneCmd struct {
	cmd *cobra.Command
}

func (f *pruneCmd) run(cmd *cobra.Command, args []string) {
	repo, err := git.OpenRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(repo)
}

func init() {
	prune := &pruneCmd{}
	prune.cmd = &cobra.Command{
		Use:   "prune",
		Short: "Remove branches with no extra commits",
		Long: `Prune looks at all of your current branches and by default will ask if you would like
		to delete ones that have no additional commits compared to its parent branch.`,
		Run: prune.run,
	}

	cli.GetRoot().AddCommand(prune.cmd)
}
