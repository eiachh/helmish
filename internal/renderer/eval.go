package renderer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// TemplateData holds the data passed to templates
type TemplateData struct {
	Values interface{}
	Chart  interface{}
}

// skipList contains token values that should skip evaluation
var skipList = []string{"{{end}}", "{{ end }}"}

func skipEval(value string) bool {
	for _, skip := range skipList {
		if value == skip {
			return true
		}
	}
	return false
}

// EvalContext holds the context for evaluating expressions
type EvalContext struct {
	Values interface{}
	Chart  interface{}
}

// NewEvalContext creates a new evaluation context with the given values and chart
func NewEvalContext(values, chart interface{}) *EvalContext {
	return &EvalContext{Values: values, Chart: chart}
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

// evaluateExpression evaluates the inner expression using text/template
func (ec *EvalContext) evaluateExpression(expr string) (interface{}, error) {
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

// EvaluateTokens goes through the tokens and evaluates TokenAction tokens using the eval context
func EvaluateTokens(tokens []Token, ctx *EvalContext) ([]Token, error) {
	for i, token := range tokens {
		if token.Type == TokenAction {
			if skipEval(token.Value) {
				continue
			}
			result, err := ctx.Evaluate(token.Value)
			if err != nil {
				return nil, err
			}
			// Update the token value with the evaluated result
			tokens[i].Value = fmt.Sprintf("%v", result)
		}
	}
	return tokens, nil
}