package flow

import (
	"fmt"
	"os"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
	"github.com/pinkluz/arcanist/lib/console"
	"github.com/pinkluz/arcanist/lib/git"
)

type flowCmd struct {
	cmd *cobra.Command

	tag   string
	atype string
	zone  string
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

		out := console.DrawGraph(*graph, nil)

		fmt.Println(out)
	case 1:
		wrk, err := repo.Worktree()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = wrk.Checkout(&gogit.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/heads/" + args[0]),
			Create: false,
			Force:  false,
			Keep:   true,
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Checked out branch %s\n", args[0])
	}
}

func init() {
	create := &flowCmd{}
	create.cmd = &cobra.Command{
		Use:   "flow",
		Short: "Show your current working tree",
		Long: `List all git branches that are tracking a local branch. Branches tracking remote braches
		are ignored and don't fit in with the arc workflow.`,
		Run: create.run,
	}

	// create.cmd.Flags().StringVarP(&create.tag, "tag", "t", "", "Asset Tag")
	// create.cmd.MarkFlagRequired("tag")

	// create.cmd.Flags().StringVarP(&create.atype, "type", "T", "", "Asset Type")
	// create.cmd.MarkFlagRequired("type")

	// create.cmd.Flags().StringVarP(&create.zone, "zone", "z", "", "Assets Zone")
	// create.cmd.MarkFlagRequired("zone")

	cli.GetRoot().AddCommand(create.cmd)
}
