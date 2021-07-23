package git

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pinkluz/arcanist/lib/globals"
	"github.com/pinkluz/arcanist/lib/util"
)

// Wrapper around raw git commands. These should be replaced when a go-git version
// of them is added or found. In most cases you might be able to get the same function
// with go-git but the performance would be horrible.

type RevListOutput struct {
	InFront int
	Behind  int
}

func execNonInteractive(cmdStr []string) (string, error) {
	cmd := exec.Command("git", cmdStr...)

	if globals.GetTrace() {
		fmt.Fprintln(os.Stderr, "trace:", cmd.Path, strings.Join(cmd.Args, " "))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return "", fmt.Errorf(string(output))
	}

	return string(output), nil
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

	stdout, err := execNonInteractive([]string{
		"rev-list", "--left-right", "--count",
		fmt.Sprintf("%s...%s", upstream, current)})
	if err != nil {
		return nil, err
	}

	split := strings.Fields(stdout)

	if len(split) == 2 {
		behind, _ := strconv.Atoi(split[0])
		ahead, _ := strconv.Atoi(split[1])

		co.Behind = behind
		co.InFront = ahead
	}

	return co, nil
}

func CheckoutTrackRaw(branch string, upstream string) error {
	_, err := execNonInteractive([]string{
		"checkout", "--track", upstream, "-b", branch})
	if err != nil {
		return err
	}

	return nil
}

func CheckoutRaw(branch string) error {
	_, err := execNonInteractive([]string{"checkout", branch})
	if err != nil {
		return err
	}

	return nil
}

func PullRebase() error {
	_, err := execNonInteractive([]string{"pull", "--rebase"})
	if err != nil {
		return err
	}

	return nil
}

func AbortRebase() error {
	_, err := execNonInteractive([]string{"rebase", "--abort"})
	if err != nil {
		return err
	}

	return nil
}

func DeleteBranch(branch string) error {
	_, err := execNonInteractive([]string{"branch", "--delete", branch})
	if err != nil {
		return err
	}

	return nil
}

func SetBranchUpstream(branch string, upstream string) error {
	_, err := execNonInteractive([]string{"branch", "-u", upstream, branch})
	if err != nil {
		return err
	}

	return nil
}

// git fetch origin mschuett/anotherone-3 +refs/heads/mschuett/anotherone-3:refs/remotes/graft/origin/mschuett/anotherone-3
func FetchWithRefspec(remote string, branch string, refspec string) error {
	_, err := execNonInteractive([]string{"fetch", remote, branch, refspec})
	if err != nil {
		return err
	}

	return nil
}

func MergeBase(refOne string, refTwo string) (string, error) {
	output, err := execNonInteractive([]string{"merge-base", refOne, refTwo})
	if err != nil {
		return "", err
	}

	// Windox/Nix compatible new line removal
	lines := util.SplitLines(output)

	if len(lines) < 1 {
		return "", fmt.Errorf("Merge base didn't return any output")
	}

	return lines[0], nil
}

func CherryPick(refOne string, refTwo string) error {
	_, err := execNonInteractive([]string{
		"cherry-pick", fmt.Sprintf("%s..%s", refOne, refTwo)})
	if err != nil {
		return err
	}

	return nil
}

func CherryPickAbort() error {
	_, err := execNonInteractive([]string{"cherry-pick", "--abort"})
	if err != nil {
		return err
	}

	return nil
}

func CommitAmendRaw() error {
	cmd := exec.Command("git", "commit", "--amend")

	env := os.Environ()
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("Commit amend failed")
	}

	return nil
}

func CommitRaw() error {
	cmd := exec.Command("git", "commit")

	env := os.Environ()
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("Commit failed")
	}

	return nil
}

func RevParseRaw(ref string) (string, error) {
	output, err := execNonInteractive([]string{"rev-parse", ref})
	if err != nil {
		return "", err
	}

	lines := util.SplitLines(output)
	if len(lines) < 1 {
		return "", fmt.Errorf("No output returned from rev-parse")
	}

	return lines[0], nil
}
