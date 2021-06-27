package flow

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
	"github.com/pinkluz/arcanist/lib/console"
	"github.com/pinkluz/arcanist/lib/git"
)

type flowCmd struct {
	cmd *cobra.Command
}

func (f *flowCmd) run(cmd *cobra.Command, args []string) {
	repo, err := git.OpenRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch len(args) {
	case 0:
		graph, err := git.GetLocalBranchGraph(repo)
		if err != nil {
			fmt.Println(err)
		}

		if graph != nil {
			out := console.DrawGraph(*graph, nil)
			fmt.Println(out)
		}
	case 1:
		err := git.Checkout(repo, args[0], "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Switched to branch '%s'", args[0])
		fmt.Println()
	case 2:
		err := git.Checkout(repo, args[0], args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Branch %s has been created", args[0])
		fmt.Println()
	}
}

func init() {
	flow := &flowCmd{}
	flow.cmd = &cobra.Command{
		Use:   "flow",
		Short: "Show your current working tree",
		Long: `List all git branches that are tracking a local branch. Branches tracking remote braches
		are ignored and don't fit in with the arc workflow.`,
		Run: flow.run,
	}

	cli.GetRoot().AddCommand(flow.cmd)
}
