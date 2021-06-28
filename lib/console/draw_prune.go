package console

import (
	"strings"
)

// Output warning and instructions for the prune command
func DrawPrune() string {
	return strings.Join([]string{
		"Prune is now showing you all branches that are scheduled for removal.",
		"Branches shown in " + prunableBranches.Render("THIS COLOR") + " will be removed. This only shows you what",
		"will be removed to remove them run this command with --destroy-branches",
	}, "\n")
}
