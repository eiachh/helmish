package ast

import (
	"fmt"
	"strings"

	"helmish/internal/renderer/types"
)

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
	n.Token.Value = fmt.Sprintf("%v", resultVal)
	*out = append(*out, n.Token)
	return nil
}

// IfNode represents an if node with condition, then, and else branches
type IfNode struct {
	Cond string // the condition expression, e.g. ".Values.enabled"
	Then []Node
	Else []Node
}

// Eval evaluates the if node
func (n *IfNode) Eval(ctx *types.EvalContext, out *[]types.Token) error {
	condResult, err := ctx.Evaluate("{{" + n.Cond + "}}")
	if err != nil {
		return err
	}
	condBool := types.IsTruthy(condResult)
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
			ifNode := &IfNode{Cond: condStr}
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
		default:
			// Skip other tokens for now
		}
		i++
	}
	return nodes, i
}