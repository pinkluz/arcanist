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
//
// $ git rev-list --left-right --count main...mschuett/test-2
// 17 1
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
