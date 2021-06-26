package git

type BranchNodeWrapper struct {
	RootNodes []*BranchNode

	// Easy access for whe runnning cascade
	BranchMap           map[string]*BranchNode
	LongestBranchLength int
}

func (b BranchNodeWrapper) IsEmpty() bool {
	return len(b.RootNodes) == 0
}

type BranchNode struct {
	Name       string
	Merge      string
	MergeShort string
	RemoteName string

	Hash          string
	CommitMsg     string
	CommitsAhead  int
	CommitsBehind int

	IsActiveBranch bool

	Upstream   *BranchNode
	Downstream []*BranchNode
}

// IsLocal tells you if this branch is tracking another branch that is in your current
// repo vs a branch in the repo repo.
func (b BranchNode) IsLocal() bool {
	return b.RemoteName == "."
}

// IsRoot checks if your branch is a parent node that has no
// upstream. You can have multiple root nodes in a single git repo.
func (b BranchNode) IsRoot() bool {
	return b.Upstream == nil
}

// CountDownstreams returns the number of downstreams from the given branch. For instance
// in the following graph output by flow...
//
// master
//  ├ mschuett/trash-test                                                 c9d98532 1:1 lol
//  └ mschuett/off-master                                                 539c188a 0:0 wowowowowo
//   └ mschuett/off-master-2                                              539c188a 0:0 wowowowowo
//
// CountDownstreams(master) -> 3
// CountDownstreams(mschuett/off-master) -> 1
// CountDownstreams(mschuett/trash-test) -> 0
func (b BranchNode) CountDownstreams() int {
	var r func(BranchNode) int

	r = func(n BranchNode) int {
		count := 0

		for _, node := range n.Downstream {
			count = count + 1
			count = count + r(*node)
		}

		return count
	}

	return r(b)
}
