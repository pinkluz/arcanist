package graft

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
	"github.com/pinkluz/arcanist/lib/git"
)

type graftCmd struct {
	cmd *cobra.Command

	remote          string
	localBranchName string
}

func (f *graftCmd) run(cmd *cobra.Command, args []string) {
	repo, err := git.OpenRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch len(args) {
	case 1:
		err := git.Graft(repo, args[0], f.remote, f.localBranchName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Grafted '%s' into your local repo", args[0])
		fmt.Println()
	default:
		cmd.Help()
		os.Exit(0)
	}
}

func init() {
	graft := &graftCmd{}
	graft.cmd = &cobra.Command{
		Use:   "graft <branch>",
		Short: "Place a branch from a given remote in your local working tree",
		Long:  `Place a branch in your local working copy from a given upstream so you can work off of it.`,
		Args:  cobra.MinimumNArgs(1),
		Run:   graft.run,
	}

	graft.cmd.Flags().StringVarP(&graft.remote, "remote", "r", "origin", "The remote to use when grafting")
	graft.cmd.Flags().StringVarP(&graft.localBranchName, "local-branch-name", "b", "", "The branch name to gaft onto")

	cli.GetRoot().AddCommand(graft.cmd)
}
