package main

import (
	"fmt"
	"strings"
)

// replaceBlock swaps the text between begin and end markers, keeping the marker
// lines. It errors if a marker is missing so a mangled target fails loudly.
func replaceBlock(content, begin, end, body string) (string, error) {
	i := strings.Index(content, begin)
	j := strings.Index(content, end)
	if i < 0 || j < 0 || j < i {
		return "", fmt.Errorf("markers %q..%q not found", begin, end)
	}
	i += len(begin)
	return content[:i] + "\n" + body + "\n" + content[j:], nil
}
