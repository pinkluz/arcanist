package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pinkluz/arcanist/lib/util"
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
func RevListRaw(current string, upstream string) (*RevListOutput, error) {

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

func CheckoutTrackRaw(branch string, upstream string) error {
	cmd := exec.Command("git", "checkout", "--track", upstream, "-b", branch)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

func CheckoutRaw(branch string) error {
	cmd := exec.Command("git", "checkout", branch)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

func PullRebase() error {
	cmd := exec.Command("git", "pull", "--rebase")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

func AbortRebase() error {
	cmd := exec.Command("git", "rebase", "--abort")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

func DeleteBranch(branch string) error {
	cmd := exec.Command("git", "branch", "--delete", branch)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

func SetBranchUpstream(branch string, upstream string) error {
	cmd := exec.Command("git", "branch", "-u", upstream, branch)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

// git fetch origin mschuett/anotherone-3 +refs/heads/mschuett/anotherone-3:refs/remotes/graft/origin/mschuett/anotherone-3
func FetchWithRefspec(remote string, branch string, refspec string) error {
	cmd := exec.Command("git", "fetch", remote, branch, refspec)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

func MergeBase(refOne string, refTwo string) (string, error) {
	cmd := exec.Command("git", "merge-base", refOne, refTwo)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return "", fmt.Errorf(string(output))
	}

	// Windox/Nix compatible new line removal
	lines := util.SplitLines(string(output))

	if len(lines) < 1 {
		return "", fmt.Errorf("Merge base didn't return any output")
	}

	return lines[0], nil
}

func CherryPick(refOne string, refTwo string) error {
	cmd := exec.Command("git", "cherry-pick", fmt.Sprintf("%s..%s", refOne, refTwo))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}

func CherryPickAbort() error {
	cmd := exec.Command("git", "cherry-pick", "--abort")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf(string(output))
	}

	return nil
}
