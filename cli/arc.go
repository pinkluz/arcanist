package cli

import (
	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli/shared"
)

var base *liliumCmd

type liliumCmd struct {
	cmd *cobra.Command

	config string
}

func (lc *liliumCmd) perprerun(cmd *cobra.Command, args []string) {
	shared.SetupConfig(lc.config)
}

func init() {
	base = &liliumCmd{}
	base.cmd = &cobra.Command{
		Use:   "arc",
		Short: "arc is a git workflow management tool",
		Long: `arc is a git workflow management tool that is based on the arcanist tooling that was build
by epriestly for Phabricator. Arc is a patch based workflow tool that focuses on creating a
readable commit log.`,
		PersistentPreRun: base.perprerun,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	base.cmd.PersistentFlags().StringVarP(&base.config, "config", "", "", "config file to use")
}

// GetRoot return sthe global root command
func GetRoot() *cobra.Command {
	return base.cmd
}
