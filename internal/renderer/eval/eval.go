package eval

import (
	"strings"

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
	applyWhitespaceTrimming(&result)
	cleanupEmptyTokens(&result)
	return result, nil
}

// cleanupEmptyTokens removes empty text tokens and trailing whitespace-only tokens
func cleanupEmptyTokens(tokens *[]types.Token) {
	// First pass: remove empty text tokens
	var cleaned []types.Token
	for _, tok := range *tokens {
		if tok.Type == types.TokenText && tok.Value == "" {
			continue
		}
		cleaned = append(cleaned, tok)
	}

	// Second pass: remove trailing whitespace-only tokens
	for len(cleaned) > 0 && cleaned[len(cleaned)-1].Type == types.TokenText && isWhitespaceOnly(cleaned[len(cleaned)-1].Value) {
		cleaned = cleaned[:len(cleaned)-1]
	}

	*tokens = cleaned
}

// applyWhitespaceTrimming processes TrimLeft and TrimRight flags on action tokens
// by removing adjacent whitespace from neighboring text tokens. It continues
// trimming across multiple tokens until hitting non-whitespace content.
func applyWhitespaceTrimming(tokens *[]types.Token) {
	for i := 0; i < len(*tokens); i++ {
		token := (*tokens)[i]
		if token.Type != types.TokenText {
			// Check for TrimLeft: trim trailing whitespace from preceding text tokens
			if token.TrimLeft && i > 0 {
				for j := i - 1; j >= 0; j-- {
					if (*tokens)[j].Type == types.TokenText {
						trimmed := trimTrailingWhitespace((*tokens)[j].Value)
						if trimmed == (*tokens)[j].Value {
							// Nothing was trimmed, check if this token has non-whitespace
							if !isWhitespaceOnly((*tokens)[j].Value) {
								// Hit non-whitespace content, stop
								break
							}
							// Token is whitespace-only but nothing to trim (already clean)
							// Continue to previous token
						} else {
							// Something was trimmed
							(*tokens)[j].Value = trimmed
							if !isWhitespaceOnly(trimmed) {
								// Still has non-whitespace content, stop
								break
							}
							// Trimmed to empty/whitespace-only, continue to previous token
						}
					} else {
						// Hit a non-text token, stop trimming
						break
					}
				}
			}
			// Check for TrimRight: trim leading whitespace from following text tokens
			if token.TrimRight && i+1 < len(*tokens) {
				for j := i + 1; j < len(*tokens); j++ {
					if (*tokens)[j].Type == types.TokenText {
						trimmed := trimLeadingWhitespace((*tokens)[j].Value)
						if trimmed == (*tokens)[j].Value {
							// Nothing was trimmed, check if this token has non-whitespace
							if !isWhitespaceOnly((*tokens)[j].Value) {
								// Hit non-whitespace content, stop
								break
							}
							// Token is whitespace-only but nothing to trim (already clean)
							// Continue to next token
						} else {
							// Something was trimmed
							(*tokens)[j].Value = trimmed
							if !isWhitespaceOnly(trimmed) {
								// Still has non-whitespace content, stop
								break
							}
							// Trimmed to empty/whitespace-only, continue to next token
						}
					} else {
						// Hit a non-text token, stop trimming
						break
					}
				}
			}
		}
	}
}

// trimTrailingWhitespace removes spaces, tabs, and newlines from the end of a string
func trimTrailingWhitespace(s string) string {
	return strings.TrimRight(s, " \t\n\r")
}

// trimLeadingWhitespace removes spaces, tabs, and newlines from the beginning of a string
func trimLeadingWhitespace(s string) string {
	return strings.TrimLeft(s, " \t\n\r")
}

// isWhitespaceOnly returns true if the string contains only spaces, tabs, and newlines
func isWhitespaceOnly(s string) bool {
	return strings.Trim(s, " \t\n\r") == ""
}

// NewEvalContext creates a new evaluation context with the given values and chart
func NewEvalContext(values, chart interface{}) *types.EvalContext {
	return &types.EvalContext{Values: values, Chart: chart, Root: values}
}