package ast

import (
	"fmt"
	"strings"

	"helmish/internal/renderer/types"
)

// CondTokenType represents the type of condition token
type CondTokenType int

const (
	CondExpr CondTokenType = iota
	CondAnd
	CondOr
	CondNot
)

// CondToken represents a token in a condition expression
type CondToken struct {
	Type  CondTokenType
	Value string
}

// CondNode represents a condition node that evaluates to a boolean
type CondNode struct {
	Tokens []CondToken
}

// Eval evaluates the condition and returns a boolean result
func (n *CondNode) Eval(ctx *types.EvalContext) (bool, error) {
	if len(n.Tokens) == 0 {
		return false, nil
	}
	if len(n.Tokens) == 1 && n.Tokens[0].Type == CondExpr {
		result, err := ctx.EvaluateSimple("{{" + n.Tokens[0].Value + "}}")
		if err != nil {
			return false, err
		}
		return types.IsTruthy(result), nil
	}
	// Check if prefix operator
	if n.Tokens[0].Type == CondAnd || n.Tokens[0].Type == CondOr || n.Tokens[0].Type == CondNot {
		op := n.Tokens[0]
		exprs := n.Tokens[1:]
		switch op.Type {
		case CondAnd:
			result := true
			for _, expr := range exprs {
				if expr.Type != CondExpr {
					return false, fmt.Errorf("invalid condition: expected expression after and")
				}
				val, err := ctx.EvaluateSimple("{{" + expr.Value + "}}")
				if err != nil {
					return false, err
				}
				if !types.IsTruthy(val) {
					result = false
				}
			}
			return result, nil
		case CondOr:
			result := false
			for _, expr := range exprs {
				if expr.Type != CondExpr {
					return false, fmt.Errorf("invalid condition: expected expression after or")
				}
				val, err := ctx.EvaluateSimple("{{" + expr.Value + "}}")
				if err != nil {
					return false, err
				}
				if types.IsTruthy(val) {
					result = true
				}
			}
			return result, nil
		case CondNot:
			if len(exprs) != 1 || exprs[0].Type != CondExpr {
				return false, fmt.Errorf("invalid condition: not expects one expression")
			}
			val, err := ctx.EvaluateSimple("{{" + exprs[0].Value + "}}")
			if err != nil {
				return false, err
			}
			return !types.IsTruthy(val), nil
		}
	} else {
		// Assume infix: expr op expr
		if len(n.Tokens) == 3 && n.Tokens[0].Type == CondExpr && n.Tokens[2].Type == CondExpr {
			left, err := ctx.EvaluateSimple("{{" + n.Tokens[0].Value + "}}")
			if err != nil {
				return false, err
			}
			right, err := ctx.EvaluateSimple("{{" + n.Tokens[2].Value + "}}")
			if err != nil {
				return false, err
			}
			switch n.Tokens[1].Type {
			case CondAnd:
				return types.IsTruthy(left) && types.IsTruthy(right), nil
			case CondOr:
				return types.IsTruthy(left) || types.IsTruthy(right), nil
			default:
				return false, fmt.Errorf("unsupported infix operator")
			}
		}
	}
	return false, fmt.Errorf("invalid condition structure")
}

// parseCondition parses a condition string into condition tokens
func parseCondition(cond string) []CondToken {
	parts := strings.Fields(cond)
	var tokens []CondToken
	for _, part := range parts {
		switch part {
		case "and":
			tokens = append(tokens, CondToken{Type: CondAnd})
		case "or":
			tokens = append(tokens, CondToken{Type: CondOr})
		case "not":
			tokens = append(tokens, CondToken{Type: CondNot})
		default:
			tokens = append(tokens, CondToken{Type: CondExpr, Value: part})
		}
	}
	return tokens
}

// Node represents a node in the AST
type Node interface {
	Eval(ctx *types.EvalContext, out *[]types.Token) error
}

// TextNode represents a text node with a token
type TextNode struct {
	Token types.Token
}

// Eval evaluates the text node
func (n *TextNode) Eval(ctx *types.EvalContext, out *[]types.Token) error {
	*out = append(*out, n.Token)
	return nil
}

// ActionNode represents an action node with a token
type ActionNode struct {
	Token types.Token
}

// Eval evaluates the action node
func (n *ActionNode) Eval(ctx *types.EvalContext, out *[]types.Token) error {
	resultVal, err := ctx.Evaluate(n.Token.Value)
	if err != nil {
		return err
	}
	// Create a new token with the evaluated value instead of modifying in place
	// This is important for range loops where the same action is evaluated multiple times
	evaluatedToken := types.Token{
		Type:   n.Token.Type,
		Value:  fmt.Sprintf("%v", resultVal),
		Line:   n.Token.Line,
		Indent: n.Token.Indent,
	}
	*out = append(*out, evaluatedToken)
	return nil
}

// IfNode represents an if node with condition, then, and else branches
type IfNode struct {
	Cond *CondNode
	Then []Node
	Else []Node
}

// Eval evaluates the if node
func (n *IfNode) Eval(ctx *types.EvalContext, out *[]types.Token) error {
	condBool, err := n.Cond.Eval(ctx)
	if err != nil {
		return err
	}
	if condBool {
		for _, node := range n.Then {
			err := node.Eval(ctx, out)
			if err != nil {
				return err
			}
		}
	} else {
		for _, node := range n.Else {
			err := node.Eval(ctx, out)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RangeNode represents a range node that iterates over a list or map
type RangeNode struct {
	Variable   string // The variable name for the current item (e.g., "." or "$item")
	Collection string // The expression to iterate over
	Body       []Node
}

// Eval evaluates the range node
func (n *RangeNode) Eval(ctx *types.EvalContext, out *[]types.Token) error {
	// Get the collection value (actual typed value, not string representation)
	result, err := ctx.GetValue(n.Collection)
	if err != nil {
		return err
	}

	// Handle different collection types
	switch collection := result.(type) {
	case []interface{}:
		for _, item := range collection {
			// Create a new context with the current item as the Values
			// This makes . refer to the current item in the range
			itemCtx := &types.EvalContext{
				Values: item,
				Chart:  ctx.Chart,
				Root:   ctx.Root,
			}
			// Evaluate each node in the body
			for _, node := range n.Body {
				err := node.Eval(itemCtx, out)
				if err != nil {
					return err
				}
			}
		}
	case map[string]interface{}:
		for _, value := range collection {
			// Create a new context with the current value as the Values
			itemCtx := &types.EvalContext{
				Values: value,
				Chart:  ctx.Chart,
				Root:   ctx.Root,
			}
			// Evaluate each node in the body
			for _, node := range n.Body {
				err := node.Eval(itemCtx, out)
				if err != nil {
					return err
				}
			}
		}
	default:
		// If it's not a collection we can iterate over, just evaluate the body once
		// or we could return an error - for now, let's skip
	}

	return nil
}

// WithNode represents a with node that re-scopes the context
type WithNode struct {
	Expression string // The expression to evaluate and re-scope to
	Body       []Node
	Else       []Node
}

// Eval evaluates the with node
func (n *WithNode) Eval(ctx *types.EvalContext, out *[]types.Token) error {
	// Get the value for the expression
	result, err := ctx.GetValue(n.Expression)
	if err != nil || !types.IsTruthy(result) {
		// Value doesn't exist or is falsy, execute else branch
		for _, node := range n.Else {
			if err := node.Eval(ctx, out); err != nil {
				return err
			}
		}
		return nil
	}

	// Create a new context with the value as the new scope
	withCtx := &types.EvalContext{
		Values: result,
		Chart:  ctx.Chart,
		Root:   ctx.Root,
	}
	for _, node := range n.Body {
		if err := node.Eval(withCtx, out); err != nil {
			return err
		}
	}
	return nil
}

// ParseAST parses a list of tokens into an AST
func ParseAST(tokens []types.Token) ([]Node, error) {
	nodes, _ := parseBlock(tokens, 0)
	return nodes, nil
}

// parseBlock parses a block of tokens until a terminator
func parseBlock(tokens []types.Token, start int, terminators ...types.TokenType) ([]Node, int) {
	var nodes []Node
	i := start
	for i < len(tokens) {
		for _, term := range terminators {
			if tokens[i].Type == term {
				return nodes, i
			}
		}
		switch tokens[i].Type {
		case types.TokenText:
			nodes = append(nodes, &TextNode{Token: tokens[i]})
		case types.TokenAction:
			nodes = append(nodes, &ActionNode{Token: tokens[i]})
		case types.TokenIf:
			// Recursive for nested if
			inner := strings.TrimSpace(strings.TrimSuffix(tokens[i].Value, "}}")[2:])
			condStr := strings.TrimSpace(inner[2:]) // remove "if"
			condTokens := parseCondition(condStr)
			condNode := &CondNode{Tokens: condTokens}
			ifNode := &IfNode{Cond: condNode}
			i++
			thenNodes, newI := parseBlock(tokens, i, types.TokenElse, types.TokenEnd)
			ifNode.Then = thenNodes
			i = newI
			if i < len(tokens) && tokens[i].Type == types.TokenElse {
				i++
				elseNodes, newI := parseBlock(tokens, i, types.TokenEnd)
				ifNode.Else = elseNodes
				i = newI
			}
			if i < len(tokens) && tokens[i].Type == types.TokenEnd {
				i++
			}
			nodes = append(nodes, ifNode)
			continue
		case types.TokenRange:
			// Parse range expression
			inner := strings.TrimSpace(strings.TrimSuffix(tokens[i].Value, "}}")[2:])
			rangeExpr := strings.TrimSpace(inner[5:]) // remove "range"
			rangeNode := &RangeNode{Collection: rangeExpr}
			i++
			bodyNodes, newI := parseBlock(tokens, i, types.TokenEnd)
			rangeNode.Body = bodyNodes
			i = newI
			if i < len(tokens) && tokens[i].Type == types.TokenEnd {
				i++
			}
			nodes = append(nodes, rangeNode)
			continue
		case types.TokenWith:
			// Parse with expression
			inner := strings.TrimSpace(strings.TrimSuffix(tokens[i].Value, "}}")[2:])
			withExpr := strings.TrimSpace(inner[4:]) // remove "with"
			withNode := &WithNode{Expression: withExpr}
			i++
			bodyNodes, newI := parseBlock(tokens, i, types.TokenElse, types.TokenEnd)
			withNode.Body = bodyNodes
			i = newI
			if i < len(tokens) && tokens[i].Type == types.TokenElse {
				i++
				elseNodes, newI := parseBlock(tokens, i, types.TokenEnd)
				withNode.Else = elseNodes
				i = newI
			}
			if i < len(tokens) && tokens[i].Type == types.TokenEnd {
				i++
			}
			nodes = append(nodes, withNode)
			continue
		default:
			// Skip other tokens for now
		}
		i++
	}
	return nodes, i
}