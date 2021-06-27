package console

import (
	"fmt"
	"strings"

	"github.com/pinkluz/arcanist/lib/git"
)

const (
	cross_mark = "\u00d7" // x
)

func DrawCascade(r git.RecursiveRebaseStatus, opts *DrawOpts) string {
	var ret []string

	if len(r.FailedBraches) > 0 {
		ret = append(ret,
			fmt.Sprintf("We couldn't rebase all your branches %d of %d failed",
				len(r.FailedBraches),
				len(r.SuccessfulBranches)+len(r.FailedBraches)))
	} else {
		ret = append(ret,
			fmt.Sprintf("All %d branches were rebased", len(r.SuccessfulBranches)))
	}

	for _, branch := range r.FailedBraches {
		ret = append(ret,
			fmt.Sprintf("%s "+cross_mark, branch))
	}

	return strings.Join(ret, "\n")
}
