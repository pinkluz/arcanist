package diff

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
	"github.com/pinkluz/arcanist/lib/git"
)

type diffCmd struct {
	cmd *cobra.Command
}

func (f *diffCmd) run(cmd *cobra.Command, args []string) {
	repo, err := git.OpenRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = git.Commit(repo, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	diff := &diffCmd{}
	diff.cmd = &cobra.Command{
		Use:   "diff",
		Short: "commit or amment to your local branch",
		Long: `diff by default will look at your local branch and check the number of commits that it is
ahead of its parent branch. If it has no commits on top of it's partent branch a commit will be created.
If it already has a commit it will ammend to the current one.`,
		Run: diff.run,
	}

	cli.GetRoot().AddCommand(diff.cmd)
}
