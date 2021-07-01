package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
)

type ReparentBranchesStatus struct {
}

func ReparentBranches(repo *gogit.Repository,
	preRemovedState BranchesAvailableForRemovalStatus) (*ReparentBranchesStatus, error) {

	if len(preRemovedState.BranchesForRemoval) == 0 {
		return nil, nil
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	if !ref.Name().IsBranch() {
		return nil, fmt.Errorf("You must be on a branch. Check if you are in a detatched state.")
	}

	if preRemovedState.FullGraphBeforeDelete == nil {
		return nil, fmt.Errorf("No branch graph from pre delete")
	}

	localRepoPreDelete := preRemovedState.FullGraphBeforeDelete.BranchMap
	for _, deletedBranch := range preRemovedState.BranchesForRemoval {
		rmedNode := localRepoPreDelete[deletedBranch.Name]
		newParent := getFirstLivingParent(rmedNode, preRemovedState.BranchesForRemoval)

		fmt.Println(newParent)
	}

	return nil, nil
}

// Find the first upstream branch of a given node provided it doesn't exist in the deletedBranches list
func getFirstLivingParent(node *BranchNode, deletedBranches []BranchNode) string {
	// This branch has no parent :(
	if node.Upstream == nil {
		return ""
	}

	if isDeletedBranch(node.Name, deletedBranches) {
		return getFirstLivingParent(node.Upstream, deletedBranches)
	}

	return node.Name
}

func isDeletedBranch(branch string, branches []BranchNode) bool {
	for _, b := range branches {
		if b.Name == branch {
			return true
		}
	}

	return false
}
