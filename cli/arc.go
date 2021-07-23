package cli

import (
	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli/shared"
	"github.com/pinkluz/arcanist/lib/globals"
)

var base *arcCmd

type arcCmd struct {
	cmd *cobra.Command

	config string
	trace  bool
}

func (lc *arcCmd) perprerun(cmd *cobra.Command, args []string) {
	shared.SetupConfig(lc.config)

	// Set global for tracing.
	if lc.trace {
		globals.SetTrace(lc.trace)
	}
}

func init() {
	base = &arcCmd{}
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
	base.cmd.PersistentFlags().BoolVarP(&base.trace, "trace", "", false, "output trace information")
}

// GetRoot return sthe global root command
func GetRoot() *cobra.Command {
	return base.cmd
}
