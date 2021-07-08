package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
)

func Commit(repo *gogit.Repository, amend bool) error {

	bnw, err := GetLocalBranchGraph(repo)
	if err != nil {
		return err
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	if !ref.Name().IsBranch() {
		return fmt.Errorf("You do not currently have a branch checked out")
	}

	branch := ref.Name().Short()
	node, ok := bnw.BranchMap[branch]
	if !ok {
		return fmt.Errorf("Branch %s not found in local branch map. Does it show up in arc flow?", branch)
	}

	if node.Upstream == nil {
		return fmt.Errorf("You are on a parent branch. Diff only works on child branches. Try git commit.")
	}

	return nil
}
