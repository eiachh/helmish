package tokens

import (
	"strings"

	"helmish/internal/renderer/types"
)

// Tokenize converts a DocumentBlocks's blocks into a list of tokens
func Tokenize(blocks types.DocumentBlocks) []types.Token {
	var tokens []types.Token
	for i, block := range blocks.Blocks {
		// Tokenize the raw content of all blocks
		blockTokens := tokenizeContent(block.Raw(), block.Line, block.Indent)
		tokens = append(tokens, blockTokens...)
		// Add newline between blocks (but not after the last block)
		if i < len(blocks.Blocks)-1 && len(tokens) > 0 {
			// Only add newline if the last token is a text token
			// Action tokens (including control structures) shouldn't have newlines appended
			if tokens[len(tokens)-1].Type == types.TokenText {
				tokens[len(tokens)-1].Value += "\n"
			} else {
				// Add a separate newline text token
				tokens = append(tokens, types.Token{Type: types.TokenText, Value: "\n", Line: block.Line})
			}
		}
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
			trimLeft := false
			trimRight := false
			// Check for {{- (trim whitespace to the left)
			if i < len(content) && content[i] == '-' {
				trimLeft = true
				i++
			}
			braceCount := 1
			for i < len(content) && braceCount > 0 {
				if content[i] == '\n' {
					line++
					i++
				} else if strings.HasPrefix(content[i:], "{{") {
					braceCount++
					i += 2
				} else if strings.HasPrefix(content[i:], "-}}") {
					// Check for -}} (trim whitespace to the right)
					trimRight = true
					braceCount--
					i += 3
				} else if strings.HasPrefix(content[i:], "}}") {
					braceCount--
					i += 2
				} else {
					i++
				}
			}
			action := content[start:i]
			// Strip the whitespace control markers from the action value
			// so downstream code sees clean {{...}} syntax
			if trimLeft && len(action) > 2 && action[2] == '-' {
				action = action[:2] + action[3:]
			}
			if trimRight && len(action) >= 4 && action[len(action)-3] == '-' {
				action = action[:len(action)-3] + action[len(action)-2:]
			}
			// Classify the action
			tokenType := classifyAction(action)
			tokens = append(tokens, types.Token{Type: tokenType, Value: action, Line: line, Indent: indent, TrimLeft: trimLeft, TrimRight: trimRight})
		} else {
			// Text until next {{ or newline
			start := i
			for i < len(content) && content[i] != '\n' && !strings.HasPrefix(content[i:], "{{") {
				i++
			}
			text := content[start:i]
			// Include newline as part of the text token if present
			if i < len(content) && content[i] == '\n' {
				text += "\n"
				line++
				i++
			}
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
	} else if strings.HasPrefix(inner, "range ") || inner == "range" {
		return types.TokenRange
	} else if strings.HasPrefix(inner, "with ") || inner == "with" {
		return types.TokenWith
	}
	return types.TokenAction
}