package git

import (
	"sort"

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

	// Sort the RootNodes for predictable output all of the downstreams have already been
	// sorted during thr call to composeBranchNodes and we can't finish the sorting until
	// the RootNodes have been placed in the BranchNodeWrapper.
	sort.Slice(bnw.RootNodes, func(i, j int) bool {
		return bnw.RootNodes[i].Name > bnw.RootNodes[j].Name
	})

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

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	currentBranch := ""
	if head.Name().IsBranch() {
		currentBranch = head.Name().Short()
	}

	for _, branch := range config.Branches {
		// This might need to be a bit more complex but seems to work for every use case I can
		// currently think of.
		ref, err := repo.Reference(plumbing.ReferenceName("refs/heads/"+branch.Name), true)
		if err != nil {
			return nil, err
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return nil, err
		}

		infront := 0
		behind := 0
		if branch.Merge.String() != "" {
			cherryOutput, err := RevListRaw(branch.Name, branch.Merge.String())
			if err != nil {
				return nil, err
			}

			// This is confusing to read but it's beacuse the way cherry outputs the differences between
			// brannches. When showing the difference in commits between the current and the upstream
			// you must invert.
			behind = cherryOutput.Behind
			infront = cherryOutput.InFront
		}

		branchNodes[branch.Name] = branchNodeWrapper{
			RootNode: branch.Merge.String() == "",
			Upstream: branch.Merge.Short(),
			Node: &BranchNode{
				Name:           branch.Name,
				Merge:          branch.Merge.String(),
				MergeShort:     branch.Merge.Short(),
				RemoteName:     branch.Remote,
				Hash:           ref.Hash().String(),
				CommitMsg:      commit.Message,
				CommitsAhead:   infront,
				CommitsBehind:  behind,
				IsActiveBranch: branch.Name == currentBranch,
				Upstream:       &BranchNode{},
				Downstream:     make([]*BranchNode, 0),
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

			// Sort the slice for sanity when being used later to graph to console
			sort.Slice(branchNodes[tmpBranch.Upstream].Node.Downstream, func(i, j int) bool {
				return branchNodes[tmpBranch.Upstream].Node.Downstream[i].Name >
					branchNodes[tmpBranch.Upstream].Node.Downstream[j].Name
			})
		}
	}

	return branchNodes, nil
}
