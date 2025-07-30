package web

import (
	"regexp"
	"strings"
)

var (
	whitespaceRegex    = regexp.MustCompile(`\s+`)
	htmlCommentRegex   = regexp.MustCompile(`<!--[\s\S]*?-->`)
	betweenTagsRegex   = regexp.MustCompile(`>\s+<`)
	leadingSpaceRegex  = regexp.MustCompile(`^\s+`)
	trailingSpaceRegex = regexp.MustCompile(`\s+$`)
)

func MinifyHTML(html string) string {
	result := html

	result = htmlCommentRegex.ReplaceAllString(result, "")

	result = betweenTagsRegex.ReplaceAllString(result, "><")

	lines := strings.Split(result, "\n")
	var minifiedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			minifiedLines = append(minifiedLines, trimmed)
		}
	}

	result = strings.Join(minifiedLines, "")

	result = whitespaceRegex.ReplaceAllString(result, " ")

	result = leadingSpaceRegex.ReplaceAllString(result, "")
	result = trailingSpaceRegex.ReplaceAllString(result, "")

	return result
}
