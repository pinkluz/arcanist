package git

import (
	"testing"
)

func TestMaxDepthZero(t *testing.T) {
	bnw := BranchNodeWrapper{
		RootNodes: []*BranchNode{
			{
				Name: "master",
			},
		},
	}

	maxDepth := bnw.MaxDepth()
	if maxDepth != 0 {
		t.Error("maxDepth expected 0 but returned", maxDepth)
	}
}

func TestMaxDepthOne(t *testing.T) {
	bnw := BranchNodeWrapper{
		RootNodes: []*BranchNode{
			{
				Name: "master",
				Downstream: []*BranchNode{
					{
						Name: "downstream-from-master",
					},
				},
			},
		},
	}

	maxDepth := bnw.MaxDepth()
	if maxDepth != 1 {
		t.Error("maxDepth expected 1 but returned", maxDepth)
	}
}

func TestMaxDepthMany(t *testing.T) {
	bnw := BranchNodeWrapper{
		RootNodes: []*BranchNode{
			{
				Name: "master",
				Downstream: []*BranchNode{
					{
						Name: "downstream-from-master",
						Downstream: []*BranchNode{
							{
								Name: "downstream-from-master-23",
							},
						},
					},
				},
			},
			{
				Name: "master2",
				Downstream: []*BranchNode{
					{
						Name: "downstream-from-master2",
					},
				},
			},
		},
	}

	maxDepth := bnw.MaxDepth()
	if maxDepth != 2 {
		t.Error("maxDepth expected 2 but returned", maxDepth)
	}
}
