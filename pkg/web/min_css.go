package web

import (
	"regexp"
	"strings"
)

var (
	cssCommentRegex      = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	cssWhitespaceRegex   = regexp.MustCompile(`\s+`)
	cssColonSpaceRegex   = regexp.MustCompile(`\s*:\s*`)
	cssSemicolonRegex    = regexp.MustCompile(`\s*;\s*`)
	cssBraceRegex        = regexp.MustCompile(`\s*{\s*`)
	cssCloseBraceRegex   = regexp.MustCompile(`\s*}\s*`)
	cssCommaRegex        = regexp.MustCompile(`\s*,\s*`)
	cssTrailingSemiRegex = regexp.MustCompile(`;(\s*})`)
)

func MinifyCSS(css string) string {
	result := css

	result = cssCommentRegex.ReplaceAllString(result, "")

	result = cssWhitespaceRegex.ReplaceAllString(result, " ")

	result = cssColonSpaceRegex.ReplaceAllString(result, ":")
	result = cssSemicolonRegex.ReplaceAllString(result, ";")
	result = cssBraceRegex.ReplaceAllString(result, "{")
	result = cssCloseBraceRegex.ReplaceAllString(result, "}")
	result = cssCommaRegex.ReplaceAllString(result, ",")

	result = cssTrailingSemiRegex.ReplaceAllString(result, "$1")

	result = strings.TrimSpace(result)

	return result
}
