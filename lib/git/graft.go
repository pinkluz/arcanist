package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
)

func Graft(repo *gogit.Repository, bnw *BranchNodeWrapper,
	branch string, remote string, localBranchName string) error {

	if remote == "" {
		remote = "origin"
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	if !ref.Name().IsBranch() {
		return fmt.Errorf("You do not currently have a branch checked out")
	}

	// Try to fetch the remote branch
	dest := fmt.Sprintf("refs/remotes/graft/%s/%s", remote, branch)
	refspec := fmt.Sprintf("+refs/heads/%s:%s", branch, dest)
	err = FetchWithRefspec(remote, branch, refspec)
	if err != nil {
		return fmt.Errorf("Unable to fetch remote branch: %s", err)
	}

	// Checkout a new branch
	if localBranchName == "" {
		localBranchName = fmt.Sprintf("graft/%s", branch)
	}
	err = CheckoutTrackRaw(localBranchName, ref.Name().Short())
	if err != nil {
		return fmt.Errorf("Unable to checkout branch it might already exist: %s", err)
	}

	// Apply all commits to the current branch
	mergePoint, err := MergeBase(ref.Name().Short(), dest)
	if err != nil {
		DeleteBranch(localBranchName)
		return fmt.Errorf("No common merge-point: %s", err)
	}

	err = CherryPick(mergePoint, dest)
	if err != nil {
		// Try to cleanup the working dir
		CherryPickAbort()
		DeleteBranch(localBranchName)
		return fmt.Errorf("Uable to cherry-pick we tried to abort for you: %s", err)
	}

	return nil
}
