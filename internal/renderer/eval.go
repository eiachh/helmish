package renderer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// evalState holds the evaluation state for handling if/end blocks
type evalState struct {
	stack []bool
}

// newEvalState creates a new evaluation state with initial active state
func newEvalState() *evalState {
	return &evalState{stack: []bool{true}}
}

// isActive returns true if the current state is active (not skipping)
func (es *evalState) isActive() bool {
	return es.stack[len(es.stack)-1]
}

// pushIf pushes a new state for an if block
func (es *evalState) pushIf(condition bool) {
	if es.isActive() {
		es.stack = append(es.stack, condition)
	} else {
		es.stack = append(es.stack, false)
	}
}

// pop removes the top state from the stack
func (es *evalState) pop() {
	if len(es.stack) > 1 {
		es.stack = es.stack[:len(es.stack)-1]
	}
}

// toggle toggles the top state of the stack
func (es *evalState) toggle() {
	if len(es.stack) > 0 {
		es.stack[len(es.stack)-1] = !es.stack[len(es.stack)-1]
	}
}

// isTruthy determines if a value is truthy
func isTruthy(v interface{}) bool {
	if s, ok := v.(string); ok {
		if s == "true" {
			return true
		} else if s == "false" {
			return false
		} else {
			return s != ""
		}
	}
	return false
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

// EvaluateTokens goes through the tokens and evaluates them using the eval context and state machine
func EvaluateTokens(tokens []Token, ctx *EvalContext) ([]Token, error) {
	state := newEvalState()
	var result []Token
	for _, token := range tokens {
		switch token.Type {
		case TokenIf:
			// Parse the condition
			inner := strings.TrimSpace(strings.TrimSuffix(token.Value, "}}")[2:])
			if !strings.HasPrefix(inner, "if") {
				return nil, fmt.Errorf("invalid if token: %s", token.Value)
			}
			condStr := strings.TrimSpace(inner[2:]) // remove "if"
			condResult, err := ctx.Evaluate("{{" + condStr + "}}")
			if err != nil {
				return nil, err
			}
			condBool := isTruthy(condResult)
			state.pushIf(condBool)
		case TokenElse:
			state.toggle()
		case TokenEnd:
			state.pop()
		default:
			if state.isActive() {
				if token.Type == TokenAction {
					inner := strings.TrimSpace(strings.TrimSuffix(token.Value, "}}")[2:])
					if inner == "else" {
						state.toggle()
						continue
					}
					resultVal, err := ctx.Evaluate(token.Value)
					if err != nil {
						return nil, err
					}
					token.Value = fmt.Sprintf("%v", resultVal)
				}
				result = append(result, token)
			}
		}
	}
	return result, nil
}