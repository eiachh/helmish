package ast

import (
	"testing"

	"helmish/internal/renderer/types"
)

func TestParseAST(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []types.Token
		expected func([]Node) bool
	}{
		{
			name: "simple if with 2 lines",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if .Values.enabled}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 1 || ifNode.Cond.Tokens[0].Type != CondExpr || ifNode.Cond.Tokens[0].Value != ".Values.enabled" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with and both true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondAnd || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with and true false",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondAnd || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with and false true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondAnd || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with and false false",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondAnd || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with and 5 conditions all true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b .Values.c .Values.d .Values.e}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 6 || ifNode.Cond.Tokens[0].Type != CondAnd {
					return false
				}
				expected := []string{".Values.a", ".Values.b", ".Values.c", ".Values.d", ".Values.e"}
				for i, exp := range expected {
					if ifNode.Cond.Tokens[i+1].Value != exp {
						return false
					}
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with or both true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondOr || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with or true false",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondOr || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with or false true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondOr || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with or false false",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondOr || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with or 5 conditions all true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b .Values.c .Values.d .Values.e}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 6 || ifNode.Cond.Tokens[0].Type != CondOr {
					return false
				}
				expected := []string{".Values.a", ".Values.b", ".Values.c", ".Values.d", ".Values.e"}
				for i, exp := range expected {
					if ifNode.Cond.Tokens[i+1].Value != exp {
						return false
					}
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with not true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if not .Values.a}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 2 || ifNode.Cond.Tokens[0].Type != CondNot || ifNode.Cond.Tokens[1].Value != ".Values.a" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with not false",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if not .Values.a}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 2 || ifNode.Cond.Tokens[0].Type != CondNot || ifNode.Cond.Tokens[1].Value != ".Values.a" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if else with or both true",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenElse, Value: "{{else}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "  key3: val3\n", Line: 5, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Type != CondOr || ifNode.Cond.Tokens[1].Value != ".Values.a" || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 || len(ifNode.Else) != 1 {
					return false
				}
				return true
			},
		},
		{
			name: "if with or infix",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if .Values.a or .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Value != ".Values.a" || ifNode.Cond.Tokens[1].Type != CondOr || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
		{
			name: "if with and infix",
			tokens: []types.Token{
				{Type: types.TokenIf, Value: "{{if .Values.a and .Values.b}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			expected: func(nodes []Node) bool {
				if len(nodes) != 1 {
					return false
				}
				ifNode, ok := nodes[0].(*IfNode)
				if !ok {
					return false
				}
				if len(ifNode.Cond.Tokens) != 3 || ifNode.Cond.Tokens[0].Value != ".Values.a" || ifNode.Cond.Tokens[1].Type != CondAnd || ifNode.Cond.Tokens[2].Value != ".Values.b" {
					return false
				}
				if len(ifNode.Then) != 2 {
					return false
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			if !tt.expected(nodes) {
				t.Errorf("Parsed nodes did not match expected structure")
			}
		})
	}
}