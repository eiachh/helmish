package types

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// TokenType represents the type of token
type TokenType int

const (
	TokenText TokenType = iota
	TokenIf
	TokenElse
	TokenEnd
	TokenAction
)

// String returns the string representation of the token type
func (t TokenType) String() string {
	switch t {
	case TokenText:
		return "Text"
	case TokenIf:
		return "If"
	case TokenElse:
		return "Else"
	case TokenEnd:
		return "End"
	case TokenAction:
		return "Action"
	default:
		return "Unknown"
	}
}

// Token represents a single token in the template
type Token struct {
	Type  TokenType
	Value string
	Line   int
	Indent int
}

// TemplateData holds the data passed to templates
type TemplateData struct {
	Values interface{}
	Chart  interface{}
}

// EvalContext holds the context for evaluating expressions
type EvalContext struct {
	Values interface{}
	Chart  interface{}
}

// Evaluate evaluates the given expression using the context
func (ec *EvalContext) Evaluate(expr string) (interface{}, error) {
	// Strip {{ }} from the expression
	if len(expr) < 4 || !strings.HasPrefix(expr, "{{") || !strings.HasSuffix(expr, "}}") {
		// Not a valid action, return as is
		return expr, nil
	}
	inner := expr[2 : len(expr)-2]
	inner = strings.TrimSpace(inner)
	// Implement expression evaluation
	result, err := ec.evaluateExpression(inner)
	if err != nil {
		// If evaluation fails, return the inner expression as is
		return inner, nil
	}
	return result, nil
}

// EvaluateSimple evaluates a simple expression without logical operators
func (ec *EvalContext) EvaluateSimple(expr string) (interface{}, error) {
	if len(expr) < 4 || !strings.HasPrefix(expr, "{{") || !strings.HasSuffix(expr, "}}") {
		return nil, fmt.Errorf("invalid expression")
	}
	inner := expr[2 : len(expr)-2]
	return ec.evaluateSimpleExpression(inner)
}

// evaluateExpression evaluates the inner expression
func (ec *EvalContext) evaluateExpression(expr string) (interface{}, error) {
	// For now, just simple expressions
	return ec.evaluateSimpleExpression(expr)
}

// evaluateSimpleExpression evaluates a simple expression using text/template
func (ec *EvalContext) evaluateSimpleExpression(expr string) (interface{}, error) {
	// Create a template with the expression wrapped in {{ }}
	tmplStr := "{{" + expr + "}}"
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}

	data := TemplateData{
		Values: ec.Values,
		Chart:  ec.Chart,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

// IsTruthy determines if a value is truthy
func IsTruthy(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		if val == "true" {
			return true
		} else if val == "false" {
			return false
		} else {
			return val != ""
		}
	case nil:
		return false
	default:
		// For numbers, etc., consider non-zero as true, but for simplicity, true
		return true
	}
}

// BlockContent represents content that can be raw or rendered
type BlockContent interface {
	Raw() string
	Rendered() string
}

// KeyValueBlock represents a YAML key-value pair (or key only if value is empty)
type KeyValueBlock struct {
	Key   string
	Value string
}

// Raw returns the raw YAML key-value pair
func (y KeyValueBlock) Raw() string {
	if y.Value == "" {
		return y.Key
	}
	return y.Key + ": " + y.Value
}

// Rendered returns the rendered YAML key-value pair
func (y KeyValueBlock) Rendered() string {
	if y.Value == "" {
		return y.Key
	}
	if strings.Contains(y.Value, "{{") && strings.Contains(y.Value, "}}") {
		// Simulate rendering by wrapping in [rendered ...]
		return y.Key + ": [rendered " + y.Value + "]"
	}
	return y.Raw()
}

// TemplateBlock represents a Helm template block
type TemplateBlock struct {
	RawContent string
}

// Raw returns the raw template content
func (t TemplateBlock) Raw() string {
	return t.RawContent
}

// Rendered returns the rendered template content
func (t TemplateBlock) Rendered() string {
	// Simulate rendering by replacing {{ ... }} with [rendered ...]
	result := strings.ReplaceAll(t.RawContent, "{{", "[rendered ")
	result = strings.ReplaceAll(result, "}}", "]")
	return result
}

// BlockType represents the type of a block
type BlockType int

const (
	KeyValueBlockType BlockType = iota
	TemplateBlockType
)

// Block represents a single line block in a rendered template
type Block struct {
	Line    int
	Type   BlockType
	Content BlockContent
	Indent int
}

// Raw returns the raw content of the block
func (b Block) Raw() string {
	return b.Content.Raw()
}

// Rendered returns the rendered content of the block
func (b Block) Rendered() string {
	return b.Content.Rendered()
}

// GetKeyValueBlock returns the KeyValueBlock if the block is of that type
func (b Block) GetKeyValueBlock() (*KeyValueBlock, bool) {
	if b.Type == KeyValueBlockType {
		if kvb, ok := b.Content.(*KeyValueBlock); ok {
			return kvb, true
		}
	}
	return nil, false
}

// GetTemplate returns the TemplateBlock if the block is of that type
func (b Block) GetTemplate() (*TemplateBlock, bool) {
	if b.Type == TemplateBlockType {
		if tb, ok := b.Content.(*TemplateBlock); ok {
			return tb, true
		}
	}
	return nil, false
}

// DocumentBlocks represents a single YAML document with its blocks
type DocumentBlocks struct {
	Blocks []Block
}

// ValueData holds the raw content and parsed data for a values file
type ValueData struct {
	Raw    string
	Parsed interface{}
}

// Values holds the chart values files (filename -> ValueData)
type Values map[string]ValueData

// Metadata holds the chart metadata files (filename -> ValueData)
type Metadata map[string]ValueData

// YamlTemplates holds the YAML template files (filename -> content)
type YamlTemplates map[string]string

// TplFiles holds the template files (filename -> content)
type TplFiles map[string]string

// Chart represents the Helm chart data
type Chart struct {
	Path          string
	Values        Values
	Metadata      Metadata
	YamlTemplates YamlTemplates
	TplFiles      TplFiles
}

// Profile represents the profile options
type Profile struct {
	Name          string
	Capabilities  Capabilities
	// Add more fields as needed
}

// Capabilities represents Helm capabilities
type Capabilities struct {
	KubeVersion string
	APIVersions []string
}

// Options holds the options for rendering
type Options struct {
	Chart   Chart
	Profile Profile
}