package css

type Stylesheet struct {
	Rules []Rule
}

type Rule interface {
	Type() RuleType
}

type RuleType int

const (
	RuleTypeStyle RuleType = iota
	RuleTypeAtRule
	RuleTypeComment
)

type StyleRule struct {
	Selectors    []string
	Declarations []Declaration
}

func (r StyleRule) Type() RuleType { return RuleTypeStyle }

type AtRule struct {
	Name    string
	Prelude string
	Rules   []Rule
}

func (r AtRule) Type() RuleType { return RuleTypeAtRule }

type Comment struct {
	Text string
}

func (c Comment) Type() RuleType { return RuleTypeComment }

type Declaration struct {
	Property  string
	Value     string
	Important bool
}
