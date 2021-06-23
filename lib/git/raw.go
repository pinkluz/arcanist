package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Wrapper around raw git commands. These should be replaced when a go-git version
// of them is added or found. In most cases you might be able to get the same function
// with go-git but the performance would be horrible.

type RevListOutput struct {
	InFront int
	Behind  int
}

// Return the difference in commits between two branches. The following examples shows
// all commits that are different between mschuett/testing and main
// $ git cherry mschuett/testing main
// + aeebbde79edd63344ee321b0f6cf6056799c557c
// + 31f1a58c1076af278f64b6523e9e85abdb8353f5
// + 8f3dbb5de9d5eb8b47dd064f416f1a4c7e3dda6c
// - 79ecb72b8211f8458d2c35930a0cff8028c46525
// + 1d76503e4e9e6c02f8b16de4d85d00f1c26cee70
// - 0cf6f549060b215ef0790f1680312dc0a39ad58f
// + 368f0058c013c6f54d32d1abc266ca0c8ff7d3a5
func RevList(current string, upstream string) (*RevListOutput, error) {

	co := &RevListOutput{
		InFront: 0,
		Behind:  0,
	}

	cmd := exec.Command("git", "rev-list", "--left-right", "--count",
		fmt.Sprintf("%s...%s", upstream, current))

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	split := strings.Fields(string(stdout))

	if len(split) == 2 {
		behind, _ := strconv.Atoi(split[0])
		ahead, _ := strconv.Atoi(split[1])

		co.Behind = behind
		co.InFront = ahead
	}

	return co, nil
}
