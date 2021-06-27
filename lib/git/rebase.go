package git

import (
	"fmt"

	pb "github.com/cheggaaa/pb/v3"
	gogit "github.com/go-git/go-git/v5"
)

type RecursiveRebaseStatus struct {
	FailedBraches      []string
	SuccessfulBranches []string
}

func RecursiveRebase(repo *gogit.Repository) (*RecursiveRebaseStatus, error) {
	bnw, err := GetLocalBranchGraph(repo)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	if !ref.Name().IsBranch() {
		// TODO make this more clear
		return nil, fmt.Errorf("You must be on a branch. Check if you are in a detatched state.")
	}

	val, ok := bnw.BranchMap[ref.Name().Short()]
	if !ok {
		return nil, fmt.Errorf("Unable to find branch in local branch map")
	}

	loadingBar := pb.StartNew(val.CountDownstreams())
	loadingBar.Set(pb.CleanOnFinish, true)
	statusMap := map[string][]string{}
	for _, node := range val.Downstream {
		failed, success, err := rebase(node, []string{}, []string{}, loadingBar)
		statusMap["failed"] = append(statusMap["failed"], failed...)
		statusMap["success"] = append(statusMap["success"], success...)
		if err != nil {
			return nil, err
		}
	}

	ret := &RecursiveRebaseStatus{
		FailedBraches:      statusMap["failed"],
		SuccessfulBranches: statusMap["success"],
	}

	err = CheckoutRaw(ref.Name().Short())
	if err != nil {
		return ret, err
	}

	loadingBar.Finish()
	return ret, err
}

func rebase(n *BranchNode, failed []string, success []string, l *pb.ProgressBar) ([]string, []string, error) {
	l.Increment()

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
		f, s, err := rebase(node, failed, success, l)
		if err != nil {
			return append(failed, f...), append(success, s...), err
		}

		failed = append(failed, f...)
		success = append(success, s...)
	}

	return failed, append(success, n.Name), err
}
