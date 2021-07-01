package git

import (
	gogit "github.com/go-git/go-git/v5"
)

type BranchesAvailableForRemovalStatus struct {
	BranchesForRemoval    []BranchNode
	FullGraphBeforeDelete *BranchNodeWrapper
}

func BranchesAvailableForRemoval(repo *gogit.Repository) (*BranchesAvailableForRemovalStatus, error) {
	bnw, err := GetLocalBranchGraph(repo)
	if err != nil {
		return nil, err
	}

	var toDelete []BranchNode
	for _, node := range bnw.RootNodes {
		toDelete = append(toDelete, findStale(node)...)
	}

	return &BranchesAvailableForRemovalStatus{
		BranchesForRemoval:    toDelete,
		FullGraphBeforeDelete: bnw,
	}, nil
}

func findStale(n *BranchNode) []BranchNode {
	var toDelete []BranchNode
	for _, dn := range n.Downstream {
		toDelete = append(toDelete, findStale(dn)...)
	}

	// We only delete nodes that have no commits vs their upstream, are not root nodes, and
	// are not the current branch that we are sitting on.
	if n.CommitsAhead > 0 || n.IsRoot() || n.IsActiveBranch {
		return toDelete
	}

	return append(toDelete, *n)
}

func DestroyBranches(b []BranchNode) error {
	for _, n := range b {
		err := DeleteBranch(n.Name)
		if err != nil {
			return err
		}
	}

	return nil
}
