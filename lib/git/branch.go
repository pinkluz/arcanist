package git

type BranchNodeWrapper struct {
	RootNodes []*BranchNode
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
	LinesAdded    int
	LinesRemoved  int

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
