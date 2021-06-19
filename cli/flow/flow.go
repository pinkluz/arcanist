package flow

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pinkluz/arcanist/cli"
)

type flowCmd struct {
	cmd *cobra.Command

	tag   string
	atype string
	zone  string
}

func (f *flowCmd) run(cmd *cobra.Command, args []string) {
	fmt.Println("flow")
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

	create.cmd.Flags().StringVarP(&create.tag, "tag", "t", "", "Asset Tag")
	create.cmd.MarkFlagRequired("tag")

	create.cmd.Flags().StringVarP(&create.atype, "type", "T", "", "Asset Type")
	create.cmd.MarkFlagRequired("type")

	create.cmd.Flags().StringVarP(&create.zone, "zone", "z", "", "Assets Zone")
	create.cmd.MarkFlagRequired("zone")

	cli.GetRoot().AddCommand(create.cmd)
}
