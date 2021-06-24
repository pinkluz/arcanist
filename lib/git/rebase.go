package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
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

	for _, node := range val.Downstream {
		err := rebase(node)
		if err != nil {
			return err
		}
	}

	err = CheckoutRaw(ref.Name().Short())
	if err != nil {
		return err
	}

	return nil
}

func rebase(n *BranchNode) error {
	err := CheckoutRaw(n.Name)
	if err != nil {
		return err
	}

	err = PullRebase()
	if err != nil {
		return err
	}

	fmt.Println("Successful rebase of " + n.Name)
	return nil
}
