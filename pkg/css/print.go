package css

import (
	"fmt"
	"strings"
)

type PrintOptions struct {
	Minify       bool
	IndentSize   int
	IndentString string
}

func DefaultPrintOptions() PrintOptions {
	return PrintOptions{
		Minify:       false,
		IndentSize:   2,
		IndentString: " ",
	}
}

func MinifyPrintOptions() PrintOptions {
	return PrintOptions{
		Minify:       true,
		IndentSize:   0,
		IndentString: "",
	}
}

func (s *Stylesheet) String() string {
	return s.Print(DefaultPrintOptions())
}

func (s *Stylesheet) Minify() string {
	return s.Print(MinifyPrintOptions())
}

func (s *Stylesheet) Print(opts PrintOptions) string {
	var builder strings.Builder

	for i, rule := range s.Rules {
		if opts.Minify && rule.Type() == RuleTypeComment {
			continue
		}

		if i > 0 && !opts.Minify {
			builder.WriteString("\n")
		}

		builder.WriteString(printRule(rule, 0, opts))

		if !opts.Minify {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func printRule(rule Rule, depth int, opts PrintOptions) string {
	switch r := rule.(type) {
	case StyleRule:
		return printStyleRule(r, depth, opts)
	case AtRule:
		return printAtRule(r, depth, opts)
	case Comment:
		return printComment(r, depth, opts)
	default:
		return ""
	}
}

func printStyleRule(rule StyleRule, depth int, opts PrintOptions) string {
	var builder strings.Builder

	indent := getIndent(depth, opts)

	builder.WriteString(indent)
	builder.WriteString(strings.Join(rule.Selectors, ", "))

	if opts.Minify {
		builder.WriteString("{")
	} else {
		builder.WriteString(" {\n")
	}

	for i, decl := range rule.Declarations {
		if !opts.Minify {
			builder.WriteString(getIndent(depth+1, opts))
		}

		builder.WriteString(decl.Property)

		if opts.Minify {
			builder.WriteString(":")
		} else {
			builder.WriteString(": ")
		}

		builder.WriteString(decl.Value)

		if decl.Important {
			if opts.Minify {
				builder.WriteString("!important")
			} else {
				builder.WriteString(" !important")
			}
		}

		builder.WriteString(";")

		if !opts.Minify && i < len(rule.Declarations)-1 {
			builder.WriteString("\n")
		}
	}

	if opts.Minify {
		builder.WriteString("}")
	} else {
		builder.WriteString("\n")
		builder.WriteString(indent)
		builder.WriteString("}")
	}

	return builder.String()
}

func printAtRule(rule AtRule, depth int, opts PrintOptions) string {
	var builder strings.Builder

	indent := getIndent(depth, opts)

	builder.WriteString(indent)
	builder.WriteString("@")
	builder.WriteString(rule.Name)

	if rule.Prelude != "" {
		if opts.Minify {
			builder.WriteString(" ")
		} else {
			builder.WriteString(" ")
		}
		builder.WriteString(rule.Prelude)
	}

	if len(rule.Rules) == 0 {
		builder.WriteString(";")
		return builder.String()
	}

	if opts.Minify {
		builder.WriteString("{")
	} else {
		builder.WriteString(" {\n")
	}

	for i, subRule := range rule.Rules {
		if opts.Minify && subRule.Type() == RuleTypeComment {
			continue
		}

		if !opts.Minify && i > 0 {
			builder.WriteString("\n")
		}

		builder.WriteString(printRule(subRule, depth+1, opts))

		if !opts.Minify {
			builder.WriteString("\n")
		}
	}

	if opts.Minify {
		builder.WriteString("}")
	} else {
		builder.WriteString(indent)
		builder.WriteString("}")
	}

	return builder.String()
}

func printComment(comment Comment, depth int, opts PrintOptions) string {
	if opts.Minify {
		return ""
	}

	indent := getIndent(depth, opts)
	return fmt.Sprintf("%s/* %s */", indent, comment.Text)
}

func getIndent(depth int, opts PrintOptions) string {
	if opts.Minify {
		return ""
	}

	totalIndent := depth * opts.IndentSize
	return strings.Repeat(opts.IndentString, totalIndent)
}
