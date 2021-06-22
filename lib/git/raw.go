package git

import (
	"os/exec"

	"github.com/pinkluz/arcanist/lib/util"
)

// Wrapper around raw git commands. These should be replaced when a go-git version
// of them is added or found. In most cases you might be able to get the same function
// with go-git but the performance would be horrible.

type CherryOutput struct {
	AdditionalCommits []string
	MissingCommits    []string

	InFront int
	Behind  int
}

// Return the difference in commits between two branches. The following examples shows
// all commits that are different between mschuett/testing and main
// $ git cherry mschuett/testing main
// + aeebbde79edd63344ee321b0f6cf6056799c557c
// + 31f1a58c1076af278f64b6523e9e85abdb8353f5
// + 8f3dbb5de9d5eb8b47dd064f416f1a4c7e3dda6c
// + 79ecb72b8211f8458d2c35930a0cff8028c46525
// + 1d76503e4e9e6c02f8b16de4d85d00f1c26cee70
// + 0cf6f549060b215ef0790f1680312dc0a39ad58f
// + 368f0058c013c6f54d32d1abc266ca0c8ff7d3a5
func Cherry(current string, upstream string) (*CherryOutput, error) {

	co := &CherryOutput{
		AdditionalCommits: []string{},
		MissingCommits:    []string{},
		InFront:           0,
		Behind:            0,
	}

	cmd := exec.Command("git", "cherry", current, upstream)

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := util.SplitLines(string(stdout))

	for _, line := range lines {
		if len(line) > 1 {
			if line[0] == byte('+') {
				co.AdditionalCommits = append(co.AdditionalCommits, line[1:])
			}

			if line[0] == byte('-') {
				co.MissingCommits = append(co.MissingCommits, line[1:])
			}
		}
	}

	co.InFront = len(co.AdditionalCommits)
	co.Behind = len(co.MissingCommits)

	return co, nil
}
