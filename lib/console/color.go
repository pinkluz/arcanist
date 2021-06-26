package console

// Feel free to add your own color profile in here and put up a PR for it. You can copy paste
// from the nocolor function to get started.

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

func nocolor(s []string, branchPadding int, activeBranch bool, root bool) []string {
	marker := ""
	if activeBranch {
		marker = current_branch
	}

	if root {
		return []string{"%s " + marker}
	}

	f := []string{
		s[0],
		s[1], // graphLine
		s[2] + marker + strings.Repeat(" ", branchPadding), // n.Name
		s[3], // hashRef
		s[4], // n.CommitsBehind
		s[5],
		s[6], // n.CommitsAhead
		s[7], // commitMsg
	}

	return f
}

func basic(s []string, branchPadding int, activeBranch bool, root bool) []string {
	if root {
		return []string{"%s %s"}
	}

	marker := ""
	if activeBranch {
		marker = current_branch
	}

	f := []string{
		color.HiBlackString(s[0]),
		color.HiBlackString(s[1]), // graphLine
		color.HiBlueString(s[2]) + color.GreenString(marker) + strings.Repeat(" ", branchPadding), // n.Name
		color.YellowString(s[3]), // hashRef
		color.RedString(s[4]),    // n.CommitsBehind
		s[5],
		color.GreenString(s[6]), // n.CommitsAhead
		s[7],                    // commitMsg
	}

	return f
}

func gloss(s []string, branchPadding int, activeBranch bool, root bool) []string {
	line := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#2D4059", Dark: "#70A1D7"})

	marker := ""
	if activeBranch {
		marker = current_branch
	}

	if root {
		return []string{"%s " + marker}
	}

	name := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#EA5455", Dark: "#F47C7C"})

	hash := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FCA3CC", Dark: "#F7F48B"})

	commit := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#222831", Dark: "#696969"})

	behind := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FA4659", Dark: "#FA4659"})

	ahead := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#2EB872", Dark: "#2EB872"})

	f := []string{
		line.Render(s[0]),
		line.Render(s[1]), // graphLine
		name.Render(s[2]) + marker + strings.Repeat(" ", branchPadding), // n.Name
		hash.Render(s[3]),   // hashRef
		behind.Render(s[4]), // n.CommitsBehind
		s[5],
		ahead.Render(s[6]),  // n.CommitsAhead
		commit.Render(s[7]), // commitMsg
	}

	return f
}
