package util

import (
	"bufio"
	"strings"
)

// SplitLines will split lines in a way that will work on strings that are formed with
// \r\n as well as \n....
func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
