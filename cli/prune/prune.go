package prune

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
	"github.com/pinkluz/arcanist/lib/console"
	"github.com/pinkluz/arcanist/lib/git"
)

type pruneCmd struct {
	cmd *cobra.Command

	destructive bool
}

func (f *pruneCmd) run(cmd *cobra.Command, args []string) {
	repo, err := git.OpenRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	status, err := git.BranchesAvailableForRemoval(repo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !f.destructive {
		graph, err := git.GetLocalBranchGraph(repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		prompt := console.DrawPrune()
		fmt.Println()
		fmt.Println(prompt)
		fmt.Println()

		out := console.DrawGraph(*graph, &console.DrawOpts{
			AvailableForDelete: status.BranchesForRemoval,
		})

		fmt.Println(out)
	} else {
		err = git.DestroyBranches(status.BranchesForRemoval)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err := git.ReparentBranches(repo, *status)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
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

	prune.cmd.Flags().BoolVarP(&prune.destructive, "destroy-branches", "d", false, "Delete branches from your local git repo")

	cli.GetRoot().AddCommand(prune.cmd)
}
