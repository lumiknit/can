package css

import (
	"strings"
	"testing"
)

func TestParseBasicCSS(t *testing.T) {
	css := `
		body {
			margin: 0;
			padding: 10px;
			color: red !important;
		}

		.btn {
			background: blue;
		}
	`

	stylesheet, err := Parse(css)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stylesheet.Rules) != 2 {
		t.Fatalf("Expected 2 rules, got %d", len(stylesheet.Rules))
	}

	// Test first rule (body)
	bodyRule, ok := stylesheet.Rules[0].(StyleRule)
	if !ok {
		t.Fatalf("Expected StyleRule, got %T", stylesheet.Rules[0])
	}

	if len(bodyRule.Selectors) != 1 || bodyRule.Selectors[0] != "body" {
		t.Errorf("Expected selector 'body', got %v", bodyRule.Selectors)
	}

	if len(bodyRule.Declarations) != 3 {
		t.Fatalf("Expected 3 declarations, got %d", len(bodyRule.Declarations))
	}

	// Test important flag
	colorDecl := bodyRule.Declarations[2]
	if colorDecl.Property != "color" || colorDecl.Value != "red" || !colorDecl.Important {
		t.Errorf("Expected 'color: red !important', got '%s: %s important=%v'",
			colorDecl.Property, colorDecl.Value, colorDecl.Important)
	}
}

func TestParseAtRule(t *testing.T) {
	css := `
		@media (max-width: 768px) {
			body {
				font-size: 14px;
			}
		}

		@import url("styles.css");
	`

	stylesheet, err := Parse(css)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stylesheet.Rules) != 2 {
		t.Fatalf("Expected 2 rules, got %d", len(stylesheet.Rules))
	}

	// Test media query
	mediaRule, ok := stylesheet.Rules[0].(AtRule)
	if !ok {
		t.Fatalf("Expected AtRule, got %T", stylesheet.Rules[0])
	}

	if mediaRule.Name != "media" {
		t.Errorf("Expected 'media', got '%s'", mediaRule.Name)
	}

	if !strings.Contains(mediaRule.Prelude, "max-width: 768px") {
		t.Errorf("Expected prelude to contain 'max-width: 768px', got '%s'", mediaRule.Prelude)
	}

	if len(mediaRule.Rules) != 1 {
		t.Fatalf("Expected 1 nested rule, got %d", len(mediaRule.Rules))
	}

	// Test import rule
	importRule, ok := stylesheet.Rules[1].(AtRule)
	if !ok {
		t.Fatalf("Expected AtRule, got %T", stylesheet.Rules[1])
	}

	if importRule.Name != "import" {
		t.Errorf("Expected 'import', got '%s'", importRule.Name)
	}
}

func TestParseComments(t *testing.T) {
	css := `
		/* This is a comment */
		body {
			color: black;
		}
		/* Another comment */
	`

	stylesheet, err := Parse(css)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stylesheet.Rules) != 3 {
		t.Fatalf("Expected 3 rules, got %d", len(stylesheet.Rules))
	}

	// Test first comment
	comment1, ok := stylesheet.Rules[0].(Comment)
	if !ok {
		t.Fatalf("Expected Comment, got %T", stylesheet.Rules[0])
	}

	if comment1.Text != "This is a comment" {
		t.Errorf("Expected 'This is a comment', got '%s'", comment1.Text)
	}
}

func TestPrintMinified(t *testing.T) {
	css := `
		/* Comment */
		body {
		/* This is a comment */
			margin: 0;
			padding: 10px;
		}

		.btn {
			background: blue;
		}
	`

	stylesheet, err := Parse(css)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	t.Logf("Original CSS:\n%s", stylesheet.String())

	minified := stylesheet.Minify()

	t.Logf("Minified CSS:\n%s", minified)

	// Should not contain comments
	if strings.Contains(minified, "/*") {
		t.Error("Minified CSS should not contain comments")
	}

	// Should not contain extra whitespace
	if strings.Contains(minified, "\n") {
		t.Error("Minified CSS should not contain newlines")
	}

	// Should contain the rules
	if !strings.Contains(minified, "body{") {
		t.Error("Minified CSS should contain 'body{'")
	}

	if !strings.Contains(minified, "margin:0;") {
		t.Error("Minified CSS should contain 'margin:0;'")
	}
}

func TestPrintPretty(t *testing.T) {
	css := `body{margin:0;padding:10px;}.btn{background:blue;}`

	stylesheet, err := Parse(css)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	pretty := stylesheet.String()

	// Should contain newlines and proper spacing
	if !strings.Contains(pretty, "body {\n") {
		t.Error("Pretty CSS should contain 'body {\\n'")
	}

	if !strings.Contains(pretty, "  margin: 0;") {
		t.Error("Pretty CSS should contain indented declarations")
	}

	if !strings.Contains(pretty, "}\n") {
		t.Error("Pretty CSS should contain closing braces with newlines")
	}
}

func TestRoundTrip(t *testing.T) {
	originalCSS := `body {
  margin: 0;
  padding: 10px;
  color: red !important;
}

.btn {
  background: blue;
  border: 1px solid black;
}`

	// Parse and print back
	stylesheet, err := Parse(originalCSS)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	printed := stylesheet.String()

	// Parse again to check consistency
	stylesheet2, err := Parse(printed)
	if err != nil {
		t.Fatalf("Second parse failed: %v", err)
	}

	if len(stylesheet.Rules) != len(stylesheet2.Rules) {
		t.Errorf("Rule count mismatch: %d vs %d", len(stylesheet.Rules), len(stylesheet2.Rules))
	}

	// Check first rule declarations
	rule1 := stylesheet.Rules[0].(StyleRule)
	rule2 := stylesheet2.Rules[0].(StyleRule)

	if len(rule1.Declarations) != len(rule2.Declarations) {
		t.Errorf("Declaration count mismatch: %d vs %d",
			len(rule1.Declarations), len(rule2.Declarations))
	}
}
