package tokens

import (
	"strings"

	"helmish/internal/renderer/types"
)

// Tokenize converts a DocumentBlocks's blocks into a list of tokens
func Tokenize(blocks types.DocumentBlocks) []types.Token {
	var tokens []types.Token
	for _, block := range blocks.Blocks {
		// Tokenize the raw content of all blocks
		blockTokens := tokenizeContent(block.Raw(), block.Line, block.Indent)
		tokens = append(tokens, blockTokens...)
	}
	return tokens
}

// tokenizeContent tokenizes the content string into Text and Action tokens
func tokenizeContent(content string, startLine int, indent int) []types.Token {
	var tokens []types.Token
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
			tokens = append(tokens, types.Token{Type: tokenType, Value: action, Line: line, Indent: indent})
		} else {
			// Text until next {{ or newline
			start := i
			for i < len(content) && content[i] != '\n' && !strings.HasPrefix(content[i:], "{{") {
				i++
			}
			text := content[start:i]
			if text != "" {
				tokens = append(tokens, types.Token{Type: types.TokenText, Value: text, Line: line, Indent: indent})
			}
		}
	}
	return tokens
}

// classifyAction determines the token type based on the action content
func classifyAction(action string) types.TokenType {
	if len(action) < 4 || !strings.HasPrefix(action, "{{") || !strings.HasSuffix(action, "}}") {
		return types.TokenAction
	}
	inner := strings.TrimSpace(strings.TrimSuffix(action, "}}")[2:])
	if strings.HasPrefix(inner, "if ") || inner == "if" {
		return types.TokenIf
	} else if inner == "else" {
		return types.TokenElse
	} else if inner == "end" {
		return types.TokenEnd
	}
	return types.TokenAction
}