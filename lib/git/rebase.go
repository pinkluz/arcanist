package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
)

const (
	failed_marker  = ""
	success_marker = ""
)

func RecursiveRebase(repo *gogit.Repository) error {
	bnw, err := GetLocalBranchGraph(repo)
	if err != nil {
		return err
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	if !ref.Name().IsBranch() {
		// TODO make this more clear
		return fmt.Errorf("You must be on a branch. Check if you are in a detatched state.")
	}

	val, ok := bnw.BranchMap[ref.Name().Short()]
	if !ok {
		return fmt.Errorf("Unable to find branch in local branch map")
	}

	statusMap := map[string][]string{}
	for _, node := range val.Downstream {
		failed, success, err := rebase(node, []string{}, []string{})
		statusMap["failed"] = append(statusMap["failed"], failed...)
		statusMap["success"] = append(statusMap["success"], success...)
		if err != nil {
			return err
		}
	}

	fmt.Println(statusMap)

	err = CheckoutRaw(ref.Name().Short())
	if err != nil {
		return err
	}

	return nil
}

func rebase(n *BranchNode, failed []string, success []string) ([]string, []string, error) {
	err := CheckoutRaw(n.Name)
	if err != nil {
		return append(failed, n.Name), success, err
	}

	err = PullRebase()
	if err != nil {
		// Before we fail try to recover from the interactive rebase.
		err := AbortRebase()
		if err != nil {
			return append(failed, n.Name), success, err
		}

		// We keep going but we have no reason to try and rebase the downstream branches
		// as they will all fail as well.
		return append(failed, n.Name), success, err
	}

	for _, node := range n.Downstream {
		f, s, err := rebase(node, failed, success)
		if err != nil {
			return append(failed, f...), append(success, s...), err
		}
	}

	return failed, append(success, n.Name), err
}
