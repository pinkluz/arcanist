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
	for _, node := range val.Downstream {
		err := rebase(node, loadingBar)
		if err != nil {
			return nil, err
		}
	}

	ret := &RecursiveRebaseStatus{
		FailedBraches:      rebaseFailedBranches,
		SuccessfulBranches: rebaseSuccessBranches,
	}

	err = CheckoutRaw(ref.Name().Short())
	if err != nil {
		return ret, err
	}

	loadingBar.Finish()
	return ret, err
}

var (
	rebaseSuccessBranches []string
	rebaseFailedBranches  []string
)

func rebase(n *BranchNode, l *pb.ProgressBar) error {
	l.Increment()

	err := CheckoutRaw(n.Name)
	if err != nil {
		rebaseFailedBranches = append(rebaseFailedBranches, n.Name)
		return err
	}

	err = PullRebase()
	if err != nil {
		// Before we fail try to recover from the interactive rebase.
		err := AbortRebase()
		if err != nil {
			rebaseFailedBranches = append(rebaseFailedBranches, n.Name)
			return err
		}

		// We keep going but we have no reason to try and rebase the downstream branches
		// as they will all fail as well.
		rebaseFailedBranches = append(rebaseFailedBranches, n.Name)
		return err
	}

	for _, node := range n.Downstream {
		err := rebase(node, l)
		if err != nil {
			rebaseFailedBranches = append(rebaseFailedBranches, n.Name)
			return err
		}
	}

	rebaseSuccessBranches = append(rebaseSuccessBranches, n.Name)
	return err
}
