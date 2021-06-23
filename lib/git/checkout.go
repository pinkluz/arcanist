package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
)

// Checkout can take three different paths depending on the inputs.
//
// checkout is set but upstream is not
// -> First check if the branch exists and check it out
// -> If the branch doesn't exist create it and set upstream to the
//    brach that we are currently on
//
// checkout is set and upstream is set
// -> Check if branch exists and if so do nothing
// -> Branch doesn't exist so we create it and set the upstream
//    to the one that has been provided.
//
// In the future this might be changed to also re-parent branches
func Checkout(repo *gogit.Repository, checkout string, upstream string) error {
	_, err := repo.Branch(checkout)
	switch err {
	case nil:
		err := CheckoutRaw(checkout)
		if err != nil {
			return err
		}
	case gogit.ErrBranchNotFound:
		// The branch wasn't found so we will try and make it
		if upstream == "" {
			ref, err := repo.Head()
			if err != nil {
				return err
			}

			if !ref.Name().IsBranch() {
				// TODO make this more clear
				return fmt.Errorf("You do not currently have a branch checked out")
			}

			err = CheckoutTrackRaw(checkout, ref.Name().Short())
			if err != nil {
				return err
			}
		} else {
			err = CheckoutTrackRaw(checkout, upstream)
			if err != nil {
				return err
			}
		}

		return nil
	default:
		return err
	}

	return nil
}
