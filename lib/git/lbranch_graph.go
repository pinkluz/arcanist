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
	branches, err := composeBranchNodes(repo)
	if err != nil {
		return nil, fmt.Errorf("unable to composeBranchNodes: %w", err)
	}

	bnw := &BranchNodeWrapper{
		RootNodes: []*BranchNode{},
		BranchMap: map[string]*BranchNode{},
	}

	for _, branch := range branches {
		if len(branch.node.Name) > bnw.LongestBranchLength {
			bnw.LongestBranchLength = len(branch.node.Name)
		}

		bnw.BranchMap[branch.node.Name] = branch.node

		// Only surface root nodes that some other branch hangs off of. A root
		// with no downstream isn't part of any stack and would just be noise.
		if branch.upstream == "" && len(branch.node.Downstream) > 0 {
			bnw.RootNodes = append(bnw.RootNodes, branch.node)
		}
	}

	// Downstreams are already sorted by composeBranchNodes; the root nodes
	// can only be sorted here once they have all been collected.
	sort.Slice(bnw.RootNodes, func(i, j int) bool {
		return bnw.RootNodes[i].Name > bnw.RootNodes[j].Name
	})

	return bnw, nil
}

// localBranch pairs a BranchNode with the name of the local branch it tracks
// ("" for root branches) until the graph pointers are wired up.
type localBranch struct {
	upstream string
	node     *BranchNode
}

// Build a tree of all branches and how they connect to each other.
func composeBranchNodes(repo *gogit.Repository) (map[string]localBranch, error) {
	config, err := repo.Config()
	if err != nil {
		return nil, fmt.Errorf("unable to get repo config: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("unable to get repo head: %w", err)
	}

	currentBranch := ""
	if head.Name().IsBranch() {
		currentBranch = head.Name().Short()
	}

	branchNodes := map[string]localBranch{}

	// Branches with tracking information configured. Only these can hang off
	// another branch in the graph.
	for _, branch := range config.Branches {
		hash, err := RevParseRaw(branch.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to get branch reference for %s: %w", branch.Name, err)
		}

		commit, err := repo.CommitObject(plumbing.NewHash(hash))
		if err != nil {
			return nil, fmt.Errorf("unable to get commit for %s: %w", branch.Name, err)
		}

		// A branch is only part of a stack when it tracks another local
		// branch (remote "."). No upstream at all, or an upstream on a real
		// remote, makes it a root of the graph.
		upstream := ""
		if branch.Merge != "" && branch.Remote == "." {
			upstream = branch.Merge.Short()
		}

		ahead, behind := 0, 0
		// The merge point can be gone if a branch with dependents was deleted
		// but never cleaned up. Skip the rev-list in that case since it would
		// throw an error.
		if branch.Merge != "" {
			if _, err := repo.Reference(branch.Merge, true); err == nil {
				revList, err := RevListRaw(branch.Name, branch.Merge.String())
				if err != nil {
					return nil, fmt.Errorf("unable to get rev-list for %s: %w", branch.Name, err)
				}

				ahead = revList.InFront
				behind = revList.Behind
			}
		}

		branchNodes[branch.Name] = localBranch{
			upstream: upstream,
			node: &BranchNode{
				Name:           branch.Name,
				Merge:          branch.Merge.String(),
				MergeShort:     branch.Merge.Short(),
				RemoteName:     branch.Remote,
				Hash:           hash,
				CommitMsg:      commit.Message,
				CommitsAhead:   ahead,
				CommitsBehind:  behind,
				IsActiveBranch: branch.Name == currentBranch,
				Downstream:     []*BranchNode{},
			},
		}
	}

	// Not all branches in your repo have tracking config. Walk all known
	// branches and add the missing ones as roots.
	branches, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("unable to get branches: %w", err)
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		if _, ok := branchNodes[name]; ok {
			return nil
		}

		branchNodes[name] = localBranch{
			node: &BranchNode{
				Name:           name,
				IsActiveBranch: name == currentBranch,
				Downstream:     []*BranchNode{},
			},
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to walk branches: %w", err)
	}

	// Now that every branch has a node, wire up the upstream and downstream
	// pointers. A branch whose upstream no longer exists (deleted but not
	// cleaned up) is left unlinked.
	for _, branch := range branchNodes {
		if branch.upstream == "" {
			continue
		}

		parent, ok := branchNodes[branch.upstream]
		if !ok {
			continue
		}

		branch.node.Upstream = parent.node
		parent.node.Downstream = append(parent.node.Downstream, branch.node)
	}

	// Sort every downstream list for predictable graph output.
	for _, branch := range branchNodes {
		downstream := branch.node.Downstream
		sort.Slice(downstream, func(i, j int) bool {
			return downstream[i].Name > downstream[j].Name
		})
	}

	return branchNodes, nil
}
