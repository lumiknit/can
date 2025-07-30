package css

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	commentRegex  = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	selectorRegex = regexp.MustCompile(`([^{]+)\{`)
	propertyRegex = regexp.MustCompile(`([^:]+):([^;]+);?`)
	atRuleRegex   = regexp.MustCompile(`@([a-zA-Z-]+)\s*([^{]*)\s*\{`)
)

func removeComments(css string) string {
	return commentRegex.ReplaceAllString(css, "")
}

func Parse(css string) (*Stylesheet, error) {
	stylesheet := &Stylesheet{}

	css = strings.TrimSpace(css)
	if css == "" {
		return stylesheet, nil
	}

	rules, err := parseRules(css)
	if err != nil {
		return nil, err
	}

	stylesheet.Rules = rules
	return stylesheet, nil
}

func parseRules(css string) ([]Rule, error) {
	var rules []Rule
	i := 0

	for i < len(css) {
		i = skipWhitespace(css, i)
		if i >= len(css) {
			break
		}

		if strings.HasPrefix(css[i:], "/*") {
			rule, newI, err := parseComment(css, i)
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)
			i = newI
		} else if css[i] == '@' {
			rule, newI, err := parseAtRule(css, i)
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)
			i = newI
		} else {
			rule, newI, err := parseStyleRule(css, i)
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)
			i = newI
		}
	}

	return rules, nil
}

func parseComment(css string, start int) (Comment, int, error) {
	end := strings.Index(css[start:], "*/")
	if end == -1 {
		return Comment{}, len(css), nil
	}

	end += start + 2
	text := css[start+2 : end-2]
	return Comment{Text: strings.TrimSpace(text)}, end, nil
}

func parseAtRule(css string, start int) (AtRule, int, error) {
	nameEnd := start + 1
	for nameEnd < len(css) && (unicode.IsLetter(rune(css[nameEnd])) || css[nameEnd] == '-') {
		nameEnd++
	}

	name := css[start+1 : nameEnd]

	i := skipWhitespace(css, nameEnd)

	braceStart := strings.IndexByte(css[i:], '{')
	if braceStart == -1 {
		semicolon := strings.IndexByte(css[i:], ';')
		if semicolon == -1 {
			return AtRule{Name: name}, len(css), nil
		}
		prelude := strings.TrimSpace(css[i : i+semicolon])
		return AtRule{Name: name, Prelude: prelude}, i + semicolon + 1, nil
	}

	braceStart += i
	prelude := strings.TrimSpace(css[i:braceStart])

	braceEnd := findMatchingBrace(css, braceStart)
	if braceEnd == -1 {
		return AtRule{Name: name, Prelude: prelude}, len(css), nil
	}

	var rules []Rule
	if braceEnd > braceStart+1 {
		var err error
		rules, err = parseRules(css[braceStart+1 : braceEnd])
		if err != nil {
			return AtRule{}, braceEnd + 1, err
		}
	}

	return AtRule{Name: name, Prelude: prelude, Rules: rules}, braceEnd + 1, nil
}

func parseStyleRule(css string, start int) (StyleRule, int, error) {
	braceStart := strings.IndexByte(css[start:], '{')
	if braceStart == -1 {
		return StyleRule{}, len(css), nil
	}

	braceStart += start
	selectorText := strings.TrimSpace(css[start:braceStart])

	selectors := parseSelectors(selectorText)

	braceEnd := findMatchingBrace(css, braceStart)
	if braceEnd == -1 {
		return StyleRule{Selectors: selectors}, len(css), nil
	}

	// Remove comments from the declaration block before parsing
	declarationBlock := css[braceStart+1 : braceEnd]
	declarationBlock = removeComments(declarationBlock)
	declarations := parseDeclarations(declarationBlock)

	return StyleRule{Selectors: selectors, Declarations: declarations}, braceEnd + 1, nil
}

func parseSelectors(text string) []string {
	selectors := strings.Split(text, ",")
	for i, selector := range selectors {
		selectors[i] = strings.TrimSpace(selector)
	}
	return selectors
}

func parseDeclarations(text string) []Declaration {
	var declarations []Declaration

	parts := strings.Split(text, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		colonIndex := strings.IndexByte(part, ':')
		if colonIndex == -1 {
			continue
		}

		property := strings.TrimSpace(part[:colonIndex])
		value := strings.TrimSpace(part[colonIndex+1:])

		important := false
		if strings.HasSuffix(value, "!important") {
			important = true
			value = strings.TrimSpace(strings.TrimSuffix(value, "!important"))
		}

		declarations = append(declarations, Declaration{
			Property:  property,
			Value:     value,
			Important: important,
		})
	}

	return declarations
}

func findMatchingBrace(css string, start int) int {
	depth := 0
	for i := start; i < len(css); i++ {
		switch css[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func skipWhitespace(css string, start int) int {
	for start < len(css) && unicode.IsSpace(rune(css[start])) {
		start++
	}
	return start
}
