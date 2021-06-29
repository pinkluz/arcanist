package git

import (
	"fmt"
	"sort"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// lbranch_graph.go contains functions for returning the graph of local
// branches ignoring any remote branch tracking.

func GetLocalBranchGraph(repo *gogit.Repository) (*BranchNodeWrapper, error) {
	bnw := &BranchNodeWrapper{
		RootNodes:           []*BranchNode{},
		BranchMap:           map[string]*BranchNode{},
		LongestBranchLength: 0,
	}

	branches, err := composeBranchNodes(repo)
	if err != nil {
		return nil, fmt.Errorf("Uable to composeBranchNodes: %s", err.Error())
	}

	for _, branch := range branches {
		if len(branch.Node.Name) > bnw.LongestBranchLength {
			bnw.LongestBranchLength = len(branch.Node.Name)
		}

		bnw.BranchMap[branch.Node.Name] = branch.Node
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
		return nil, fmt.Errorf("Uable to get repo config: %s", err.Error())
	}

	// Temp struct to allow us to quickly create the graph at
	// the expense of a tiny bit more memory.
	branchNodes := map[string]branchNodeWrapper{}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("Uable to get repo head: %s", err.Error())
	}

	currentBranch := ""
	if head.Name().IsBranch() {
		currentBranch = head.Name().Short()
	}

	for _, branch := range config.Branches {
		// This might need to be a bit more complex but seems to work for every use case I can
		// currently think of.
		ref, err := repo.Reference(plumbing.ReferenceName("/refs/heads/"+branch.Name), true)
		if err != nil {
			return nil, fmt.Errorf("Uable to get branch reference for %s: %s", branch.Name, err.Error())
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return nil, fmt.Errorf("Uable to get commitObject: %s", err.Error())
		}

		infront := 0
		behind := 0
		// Check if the branch merge point exists still. This happens if you delete a branch with a dependency but
		// don't clean it up. In this case we just skip trying to get the rev-list which will throw an error.
		_, err = repo.Reference(plumbing.ReferenceName(branch.Merge.String()), true)
		if branch.Merge.String() != "" && err == nil {
			cherryOutput, err := RevListRaw(branch.Name, branch.Merge.String())
			if err != nil {
				return nil, fmt.Errorf("Uable to get rev-list: %s", err.Error())
			}

			behind = cherryOutput.Behind
			infront = cherryOutput.InFront
		}

		// Special cases everywhere. Git is hard to work with.
		isRootNode := false
		switch {
		case branch.Merge.String() == "":
			isRootNode = true
		case ref.Name().String() == "/"+branch.Merge.String():
			isRootNode = true
		}

		branchUpstream := ""
		if !isRootNode {
			branchUpstream = branch.Merge.Short()
		}

		branchNodes[branch.Name] = branchNodeWrapper{
			RootNode: isRootNode,
			Upstream: branchUpstream,
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
				Upstream:       nil,
				Downstream:     make([]*BranchNode, 0),
			},
		}
	}

	// Not all branches that are tracking in your repo will be listed in your config.
	// Get all known branches to make sure we have everything.
	branches, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("Uable to get braches: %s", err.Error())
	}

	branches.ForEach(func(ref *plumbing.Reference) error {
		_, ok := branchNodes[ref.Name().Short()]
		if !ok {
			branchNodes[ref.Name().Short()] = branchNodeWrapper{
				RootNode: true,
				Upstream: "",
				Node: &BranchNode{
					Name:           ref.Name().Short(),
					Merge:          "",
					MergeShort:     "",
					RemoteName:     "",
					IsActiveBranch: ref.Name().Short() == currentBranch,
					Upstream:       nil,
					Downstream:     make([]*BranchNode, 0),
				},
			}
		}

		return nil
	})

	// Set all upstream and downstreams in the BranchNodes now that we have
	// fully populated all branches.
	for _, tmpBranch := range branchNodes {
		// if branch upstream doesn't exist just skip it. This means that the branch upstream
		// was deleted but child branches were not cleaned up.
		branchUpstream, ok := branchNodes[tmpBranch.Upstream]
		if !ok {
			continue
		}

		if !tmpBranch.RootNode {
			tmpBranch.Node.Upstream = branchUpstream.Node
		}

		if tmpBranch.Upstream != "" {
			branchUpstream.Node.Downstream =
				append(branchUpstream.Node.Downstream, tmpBranch.Node)

			// Sort the slice for sanity when being used later to graph to console
			sort.Slice(branchUpstream.Node.Downstream, func(i, j int) bool {
				return branchUpstream.Node.Downstream[i].Name >
					branchUpstream.Node.Downstream[j].Name
			})
		}
	}

	return branchNodes, nil
}
