package renderer

import "strings"

// TokenType represents the type of token
type TokenType int

const (
	TokenText TokenType = iota
	TokenAction
)

// String returns the string representation of the token type
func (t TokenType) String() string {
	switch t {
	case TokenText:
		return "Text"
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

// Tokenize converts a RenderedTemplate's blocks into a list of tokens
func Tokenize(rendered RenderedTemplate) []Token {
	var tokens []Token
	for _, block := range rendered.Blocks {
		// Tokenize the raw content of all blocks
		blockTokens := tokenizeContent(block.Raw(), block.Line)
		tokens = append(tokens, blockTokens...)
	}
	return tokens
}

// tokenizeContent tokenizes the content string into Text and Action tokens
func tokenizeContent(content string, startLine int) []Token {
	var tokens []Token
	lines := strings.Split(content, "\n")
	for i, l := range lines {
		line := startLine + i
		indent := len(l) - len(strings.TrimLeft(l, " "))
		lineTokens := tokenizeLine(l, line, indent)
		tokens = append(tokens, lineTokens...)
	}
	return tokens
}

// tokenizeLine tokenizes a single line into Text and Action tokens
func tokenizeLine(content string, line int, indent int) []Token {
	var tokens []Token
	i := 0
	for i < len(content) {
		if strings.HasPrefix(content[i:], "{{") {
			// Start of an action
			start := i
			i += 2
			braceCount := 1
			for i < len(content) && braceCount > 0 {
				if strings.HasPrefix(content[i:], "{{") {
					braceCount++
					i += 2
				} else if strings.HasPrefix(content[i:], "}}") {
					braceCount--
					i += 2
				} else {
					i++
				}
			}
			action := content[start:i]
			tokens = append(tokens, Token{Type: TokenAction, Value: action, Line: line, Indent: indent})
		} else {
			// Text until next {{
			start := i
			for i < len(content) && !strings.HasPrefix(content[i:], "{{") {
				i++
			}
			text := content[start:i]
			if text != "" {
				tokens = append(tokens, Token{Type: TokenText, Value: text, Line: line, Indent: indent})
			}
		}
	}
	return tokens
}