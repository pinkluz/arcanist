package git

import (
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// lbranch_graph.go contains functions for returning the graph of local
// branches ignoring any remote branch tracking.

func GetLocalBranchGraph(repo *gogit.Repository) (*BranchNodeWrapper, error) {
	bnw := &BranchNodeWrapper{
		RootNodes: []*BranchNode{},
	}

	branches, err := composeBranchNodes(repo)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		// Make sure it's a root node and that some other branch cares about it
		// being in the graph. If it has no downstream it isn't being used at all.
		if branch.RootNode && len(branch.Node.Downstream) > 0 {
			bnw.RootNodes = append(bnw.RootNodes, branch.Node)
		}
	}

	return bnw, nil
}

type branchNodeWrapper struct {
	RootNode bool
	Upstream string
	Node     *BranchNode
}

// Build a tree of all branches and how they connect to each other
func composeBranchNodes(repo *gogit.Repository) (map[string]branchNodeWrapper, error) {
	config, err := repo.Config()
	if err != nil {
		return nil, err
	}

	// Temp struct to allow us to quickly create the graph at
	// the expense of a tiny bit more memory.
	branchNodes := map[string]branchNodeWrapper{}

	for _, branch := range config.Branches {
		branchNodes[branch.Name] = branchNodeWrapper{
			RootNode: branch.Merge.String() == "",
			Upstream: branch.Merge.Short(),
			Node: &BranchNode{
				Name:       branch.Name,
				Merge:      branch.Merge.String(),
				MergeShort: branch.Merge.Short(),
				RemoteName: branch.Remote,
				Upstream:   &BranchNode{},
				Downstream: make([]*BranchNode, 0),
			},
		}
	}

	// Not all branches that are tracking in your repo will be listed in your config.
	// Get all known branches to make sure we have everything.
	branches, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	branches.ForEach(func(ref *plumbing.Reference) error {
		_, ok := branchNodes[ref.Name().Short()]
		if !ok {
			branchNodes[ref.Name().Short()] = branchNodeWrapper{
				RootNode: true,
				Upstream: "",
				Node: &BranchNode{
					Name:       ref.Name().Short(),
					Merge:      "",
					MergeShort: "",
					RemoteName: "",
					Upstream:   &BranchNode{},
					Downstream: make([]*BranchNode, 0),
				},
			}
		}

		return nil
	})

	// Set all upstream and downstreams in the BranchNodes now that we have
	// fully populated all branches.
	for _, tmpBranch := range branchNodes {
		if !tmpBranch.RootNode {
			tmpBranch.Node.Upstream = branchNodes[tmpBranch.Upstream].Node
		}

		if tmpBranch.Upstream != "" {
			branchNodes[tmpBranch.Upstream].Node.Downstream =
				append(branchNodes[tmpBranch.Upstream].Node.Downstream, tmpBranch.Node)
		}
	}

	return branchNodes, nil
}
