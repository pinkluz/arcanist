package console

import (
	"fmt"
	"strings"

	"github.com/pinkluz/arcanist/lib/git"
	"github.com/pinkluz/arcanist/lib/util"
)

const (
	// You can use a site like
	// https://unicode-table.com/en to find the
	// character you want.  Typically double byte characters are referred to
	// like this: "Unicode Character 'CROSS MARK' (U+274C)", resulting in a
	// golang string like this: "\u274c".

	vertical_right = "\u251c" // |-
	vertical       = "\u2502" // |
	horizontal     = "\u2500" // -
	up_right       = "\u2514" // |_

	current_branch = "\u0e4f" // o
)

// DrawGraphOpts is left for later when we want to allow the user some control
// over the output. For now it's just empty.
type DrawOpts struct {
	AvailableForDelete []git.BranchNode
}

// DrawGraph takes a git.BranchNodeWrapper and renders the output for your
// console. This is all returned as a string so you can
func DrawGraph(bnw git.BranchNodeWrapper, opts *DrawOpts) string {
	internalOpts := opts
	if opts == nil {
		internalOpts = &DrawOpts{}
	}

	// Not much to do here
	if bnw.IsEmpty() {
		return "No branches found"
	}

	maxDepth := bnw.MaxDepth()

	// This could be part of the recursive drawLine call but is pulled out to
	// reduce the complexity and allow me to understand this when I read it later.
	lines := []string{}
	for _, node := range bnw.RootNodes {
		lines = append(lines, walkNodes(internalOpts, node, 0, []int{}, false, bnw.LongestBranchLength, maxDepth)...)
	}

	output := strings.Join(lines, "\n")
	return strings.TrimSuffix(output, "\n")
}

// Render a single line of the flow output
func walkNodes(o *DrawOpts, n *git.BranchNode, depth int,
	openDepths []int, cap bool, longestBranch int, maxDepth int) []string {

	var lines []string

	if depth > 0 {
		openDepths = append(openDepths, depth)
	}

	if cap {
		openDepths = removeDepthFromSlice(openDepths, depth)
	}

	// If this is a root node we don't output much information for it. If you are using flow
	// like intended this should matter but may lead to some confusion if you start committing
	// to a root branch (likely main or master).
	lines = append(lines, drawLine(*o, n, depth, openDepths, cap, longestBranch, maxDepth))

	if len(n.Downstream) > 0 {
		for i, node := range n.Downstream {
			// determine if this is the last node at a given depth and if the |_ bracket
			// needs to be used instead of a |-.
			cap := false
			if i+1 == len(n.Downstream) {
				cap = true
			}

			lines = append(lines, walkNodes(o, node, depth+1, openDepths, cap, longestBranch, maxDepth)...)
		}
	}

	return lines
}

func removeDepthFromSlice(s []int, depth int) []int {
	index := -1
	for idx, num := range s {
		if num == depth {
			index = idx
			break
		}
	}

	ret := make([]int, 0)

	// Didn't find desired depth to remove so we return a copy of the
	// same array.
	if index == -1 {
		return append(ret, s...)
	}

	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func isDepthOpen(openDepths []int, depth int) bool {
	for _, openDepth := range openDepths {
		if openDepth == depth {
			return true
		}
	}

	return false
}

// draw a single line taking into account all of the options passed in
func drawLine(o DrawOpts, n *git.BranchNode, depth int,
	openDepths []int, cap bool, longestBranch int, maxDepth int) string {

	padding := ""

	for i := 0; i < depth; i++ {
		if i != depth && isDepthOpen(openDepths, i) {
			padding = padding + vertical
		} else {
			padding = padding + " "
		}
	}

	graphLine := vertical_right
	if cap {
		graphLine = up_right
	}

	// Edge case for when printing out a root node. We stop and return here
	// and don't do a bunch of extra formatting.
	if depth == 0 {
		rootFmt := gloss([]string{}, 0, n.IsActiveBranch, true, false)
		return fmt.Sprintf(strings.Join(rootFmt, ""), n.Name)
	}

	commitMsg := "[no commit message found]"
	commitLines := util.SplitLines(n.CommitMsg)
	if len(commitLines) > 0 {
		commitMsg = commitLines[0]
	}

	defaultBranchPadding := 40
	if defaultBranchPadding < (longestBranch + maxDepth) {
		// +3 is notable here because the padding is used with strings.Repeat which can't have
		// a negative number. We do +1 because we want a space between the branch name and whatever
		// comes after it. We do another +2 because if it's the active branch we need to leave space
		// for the current branch marker to be rendered and leave a space on both sides.
		defaultBranchPadding = longestBranch + maxDepth + 3
	}

	// Padding to add to the end of every branch name. This should be a reasonable number
	// that fits most usecases. The output should look good unless you are using super
	// long branch names like wow-i-messed-up-that-last-commit-hope-this-fixes-it
	branchPadding := defaultBranchPadding - depth - len(n.Name)
	if n.IsActiveBranch {
		branchPadding = branchPadding - 1
	}

	// If somehow I still messed up all the calculations in some edge cases make sure
	// that branchPadding is never passed as a number less than 0 or the program panics.
	if branchPadding < 0 {
		branchPadding = 0
	}

	fmtStr := []string{
		padding,
		"%s ", // graphLine
		"%s ", // n.Name
		"%s ", // hashRef
		"%d",  // n.CommitsBehind
		":",
		"%d ", // n.CommitsAhead
		"%s",  // commitMsg
	}

	scheduledForRemoval := selectedForRemoval(o.AvailableForDelete, *n)
	fmtStr = gloss(fmtStr, branchPadding, n.IsActiveBranch, false, scheduledForRemoval)

	hashRef := ""
	if len(n.Hash) >= 8 {
		hashRef = n.Hash[:8]
	}

	return fmt.Sprintf(strings.Join(fmtStr, ""),
		graphLine, n.Name, hashRef, n.CommitsBehind, n.CommitsAhead, commitMsg)
}

func selectedForRemoval(nodes []git.BranchNode, node git.BranchNode) bool {
	for _, n := range nodes {
		if n.Name == node.Name {
			return true
		}
	}

	return false
}
