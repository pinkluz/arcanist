package console

import (
	"fmt"
	"strings"

	"github.com/pinkluz/arcanist/lib/git"
)

const (
	// You can use a site like
	// https://unicode-table.com/en to find the
	// character you want.  Typically double byte characters are referred to
	// like this: "Unicode Character 'CROSS MARK' (U+274C)", resulting in a
	// two byte string like this: "\u274c".

	vertical_left  = "\u2527" // -|
	vertical_right = "\u251c" // |-
	vertical       = "\u2502" // |
	horizontal     = "\u2500" // -
	up_right       = "\u2514" // |_
)

// DrawGraphOpts is left for later when we want to allow the user some control
// over the output. For now it's just empty.
type DrawGraphOpts struct {
}

// DrawGraph takes a git.BranchNodeWrapper and renders the output for your
// console. This is all returned as a string so you can
func DrawGraph(bnw git.BranchNodeWrapper, opts *DrawGraphOpts) string {
	// internalOpts := opts
	// if opts == nil {
	// 	internalOpts = &DrawGraphOpts{}
	// }

	// Not much to do here
	if bnw.IsEmpty() {
		return ""
	}

	// This could be part of the recursive drawLine call but is pulled out to
	// reduce the complexity and allow me to understand this when I read it later.
	lines := []string{}
	for _, node := range bnw.RootNodes {
		lines = append(lines, walkNodes(node, 0, false)...)
	}

	output := strings.Join(lines, "\n")
	return strings.TrimSuffix(output, "\n")
}

// Render a single line of the flow output
func walkNodes(n *git.BranchNode, depth int, final bool) []string {
	var lines []string
	// If this is a root node we don't output much information for it. If you are using flow
	// like intended this should matter but may lead to some confusion if you start committing
	// to a root branch (likely main or master).
	lines = append(lines, drawLine(n, depth, final))

	if len(n.Downstream) > 0 {
		for i, node := range n.Downstream {
			isLastNode := false
			if len(n.Downstream) == i+1 {
				isLastNode = true
			}
			lines = append(lines, walkNodes(node, depth+1, isLastNode)...)
		}
	}

	return lines
}

// draw a single line taking into account all of the options passed in
func drawLine(n *git.BranchNode, depth int, final bool) string {
	leadingSpace := strings.Repeat(" ", depth)

	graphLine := vertical_right + horizontal
	if final {
		graphLine = up_right + horizontal
	}

	// Edge case for when printing out a root node
	if depth == 0 {
		return fmt.Sprintf("%s", n.Name)
	}

	return fmt.Sprintf(leadingSpace+"%s %s", graphLine, n.Name)
}
