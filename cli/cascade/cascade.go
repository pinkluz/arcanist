package cascade

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
	"github.com/pinkluz/arcanist/lib/console"
	"github.com/pinkluz/arcanist/lib/git"
)

type cascadeCmd struct {
	cmd *cobra.Command
}

func (f *cascadeCmd) run(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Help()
		os.Exit(1)
	}

	repo, err := git.OpenRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	status, err := git.RecursiveRebase(repo)
	if err != nil {
		fmt.Println(err)
	}

	if status != nil {
		fmt.Println(console.DrawCascade(*status, nil))
	}
}

func init() {
	cascade := &cascadeCmd{}
	cascade.cmd = &cobra.Command{
		Use:   "cascade",
		Short: "rebase all the way down the chain",
		Long: `Discovers all dependencies that you have in your local branches and recursively trys to rebase
them based on the branch you run cascade from.`,
		Run: cascade.run,
	}

	cli.GetRoot().AddCommand(cascade.cmd)
}
