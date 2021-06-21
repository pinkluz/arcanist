package console

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pinkluz/arcanist/lib/git"
)

func TestSimpleDrawGraph(t *testing.T) {
	bnw := &git.BranchNodeWrapper{
		RootNodes: []*git.BranchNode{
			{
				Name:       "main",
				Merge:      "",
				MergeShort: "",
				RemoteName: "",
				Upstream:   &git.BranchNode{},
				Downstream: []*git.BranchNode{
					{
						Name:       "main/branch-1",
						Merge:      "refs/heads/main",
						MergeShort: "main",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
				},
			},
			{
				Name:       "master",
				Merge:      "",
				MergeShort: "",
				RemoteName: "",
				Upstream:   &git.BranchNode{},
				Downstream: []*git.BranchNode{
					{
						Name:       "master/branch-1",
						Merge:      "refs/heads/master",
						MergeShort: "master",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
				},
			},
		},
	}

	out := DrawGraph(*bnw, nil)

	expected := strings.TrimSpace(`
main
 └─ main/branch-1
master
 └─ master/branch-1`)

	if out != expected {
		t.Error("TestSimpleDrawGraph failed")
		t.Log("Got:")
		fmt.Println(out)
		t.Log("Expected:")
		fmt.Println(expected)
	}
}

func TestComplexDrawGraph(t *testing.T) {
	bnw := &git.BranchNodeWrapper{
		RootNodes: []*git.BranchNode{
			{
				Name:       "main",
				Merge:      "",
				MergeShort: "",
				RemoteName: "",
				Upstream:   &git.BranchNode{},
				Downstream: []*git.BranchNode{
					{
						Name:       "main/branch-1",
						Merge:      "refs/heads/main",
						MergeShort: "main",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
					{
						Name:       "main/branch-2",
						Merge:      "refs/heads/main",
						MergeShort: "main",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{
							{
								Name:       "main/branch-2-1",
								Merge:      "refs/heads/main/branch-2",
								MergeShort: "main/branch-2",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{},
							},
							{
								Name:       "main/branch-2-2",
								Merge:      "refs/heads/main/branch-2",
								MergeShort: "main/branch-2",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{
									{
										Name:       "main/branch-2-2-1",
										Merge:      "refs/heads/main/branch-2-2",
										MergeShort: "main/branch-2-2",
										RemoteName: ".",
										Upstream:   &git.BranchNode{},
										Downstream: []*git.BranchNode{},
									},
								},
							},
							{
								Name:       "main/branch-2-3",
								Merge:      "refs/heads/main/branch-2",
								MergeShort: "main/branch-2",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{},
							},
							{
								Name:       "main/branch-2-4",
								Merge:      "refs/heads/main/branch-2",
								MergeShort: "main/branch-2",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{
									{
										Name:       "main/branch-2-4-1",
										Merge:      "refs/heads/main/branch-2-4",
										MergeShort: "main/branch-2-4",
										RemoteName: ".",
										Upstream:   &git.BranchNode{},
										Downstream: []*git.BranchNode{},
									},
								},
							},
						},
					},
					{
						Name:       "main/branch-3",
						Merge:      "refs/heads/main",
						MergeShort: "main",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
					{
						Name:       "main/branch-4",
						Merge:      "refs/heads/main",
						MergeShort: "main",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
					{
						Name:       "main/branch-5",
						Merge:      "refs/heads/main",
						MergeShort: "main",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{
							{
								Name:       "main/branch-5-1",
								Merge:      "refs/heads/main/branch-5",
								MergeShort: "main/branch-5",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{},
							},
							{
								Name:       "main/branch-5-2",
								Merge:      "refs/heads/main/branch-5",
								MergeShort: "main/branch-5",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{
									{
										Name:       "main/branch-5-2-1",
										Merge:      "refs/heads/main/branch-5-2",
										MergeShort: "main/branch-5-2",
										RemoteName: ".",
										Upstream:   &git.BranchNode{},
										Downstream: []*git.BranchNode{},
									},
								},
							},
							{
								Name:       "main/branch-5-3",
								Merge:      "refs/heads/main/branch-5",
								MergeShort: "main/branch-5",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{},
							},
							{
								Name:       "main/branch-5-4",
								Merge:      "refs/heads/main/branch-5",
								MergeShort: "main/branch-5",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{},
							},
						},
					},
				},
			},
			{
				Name:       "master",
				Merge:      "",
				MergeShort: "",
				RemoteName: "",
				Upstream:   &git.BranchNode{},
				Downstream: []*git.BranchNode{
					{
						Name:       "master/branch-1",
						Merge:      "refs/heads/master",
						MergeShort: "master",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
					{
						Name:       "master/branch-2",
						Merge:      "refs/heads/master",
						MergeShort: "master",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{
							{
								Name:       "master/branch-2-1",
								Merge:      "refs/heads/master/branch-2",
								MergeShort: "master/branch-2",
								RemoteName: ".",
								Upstream:   &git.BranchNode{},
								Downstream: []*git.BranchNode{},
							},
						},
					},
					{
						Name:       "master/branch-3",
						Merge:      "refs/heads/master",
						MergeShort: "master",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
					{
						Name:       "master/branch-4",
						Merge:      "refs/heads/master",
						MergeShort: "master",
						RemoteName: ".",
						Upstream:   &git.BranchNode{},
						Downstream: []*git.BranchNode{},
					},
				},
			},
		},
	}

	out := DrawGraph(*bnw, nil)

	expected := strings.TrimSpace(`
main
 └─ main/branch-1
master
 └─ master/branch-1`)

	if out != expected {
		t.Error("TestSimpleDrawGraph failed")
		t.Log("Got:")
		fmt.Println(out)
		t.Log("Expected:")
		fmt.Println(expected)
	}
}

func TestDrawLines(t *testing.T) {
	node := &git.BranchNode{
		Name:       "main/branch-5-2-1",
		Merge:      "refs/heads/main/branch-5-2",
		MergeShort: "main/branch-5-2",
		RemoteName: ".",
		Upstream:   &git.BranchNode{},
		Downstream: []*git.BranchNode{},
	}

	out := drawLine(node, 3, []int{2}, true)

	fmt.Println(out)
}
