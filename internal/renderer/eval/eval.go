package eval

import (
	"helmish/internal/renderer/ast"
	"helmish/internal/renderer/types"
)

// EvaluateAST evaluates the AST nodes
func EvaluateAST(nodes []ast.Node, ctx *types.EvalContext) ([]types.Token, error) {
	var result []types.Token
	for _, node := range nodes {
		err := node.Eval(ctx, &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// NewEvalContext creates a new evaluation context with the given values and chart
func NewEvalContext(values, chart interface{}) *types.EvalContext {
	return &types.EvalContext{Values: values, Chart: chart}
}