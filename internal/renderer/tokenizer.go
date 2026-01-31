package renderer

import "strings"

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

// Tokenize converts a RenderedTemplate's blocks into a list of tokens
func Tokenize(rendered RenderedTemplate) []Token {
	var tokens []Token
	for _, block := range rendered.Blocks {
		// Tokenize the raw content of all blocks
		blockTokens := tokenizeContent(block.Raw(), block.Line, block.Indent)
		tokens = append(tokens, blockTokens...)
	}
	return tokens
}

// tokenizeContent tokenizes the content string into Text and Action tokens
func tokenizeContent(content string, startLine int, indent int) []Token {
	var tokens []Token
	i := 0
	line := startLine
	for i < len(content) {
		if content[i] == '\n' {
			line++
			i++
			continue
		}
		if strings.HasPrefix(content[i:], "{{") {
			// Start of an action
			start := i
			i += 2
			braceCount := 1
			for i < len(content) && braceCount > 0 {
				if content[i] == '\n' {
					line++
					i++
				} else if strings.HasPrefix(content[i:], "{{") {
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
			// Classify the action
			tokenType := classifyAction(action)
			tokens = append(tokens, Token{Type: tokenType, Value: action, Line: line, Indent: indent})
		} else {
			// Text until next {{ or newline
			start := i
			for i < len(content) && content[i] != '\n' && !strings.HasPrefix(content[i:], "{{") {
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

// tokenizeLine tokenizes a single line into Text, If, End, and Action tokens
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
			// Classify the action
			tokenType := classifyAction(action)
			tokens = append(tokens, Token{Type: tokenType, Value: action, Line: line, Indent: indent})
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

// classifyAction determines the token type based on the action content
func classifyAction(action string) TokenType {
	if len(action) < 4 || !strings.HasPrefix(action, "{{") || !strings.HasSuffix(action, "}}") {
		return TokenAction
	}
	inner := strings.TrimSpace(strings.TrimSuffix(action, "}}")[2:])
	if strings.HasPrefix(inner, "if ") || inner == "if" {
		return TokenIf
	} else if inner == "else" {
		return TokenElse
	} else if inner == "end" {
		return TokenEnd
	}
	return TokenAction
}