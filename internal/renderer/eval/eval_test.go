package eval_test

import (
	"reflect"
	"testing"

	"helmish/internal/renderer/ast"
	"helmish/internal/renderer/eval"
	"helmish/internal/renderer/types"
)

func TestEvaluateASTIf(t *testing.T) {
	// Helper to create EvalContext with given values
	createCtx := func(values map[string]interface{}) *types.EvalContext {
		return eval.NewEvalContext(values, map[string]interface{}{"name": "test"})
	}

	tests := []struct {
		name     string
		tokens   []types.Token
		values   map[string]interface{}
		expected []types.Token
	}{
		{
			name: "just text tokens",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n  key: value\n", Line: 1, Indent: 0},
			},
			values: map[string]interface{}{},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n  key: value\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "text, if true, keyvalue inside",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.enabled}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"enabled": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "text, if false, keyvalue inside",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.enabled}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"enabled": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "text with keyvalue after false if",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.enabled}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key1: value1\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "  key2: value2\n", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{"enabled": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key2: value2\n", Line: 5, Indent: 0},
			},
		},
		{
			name: "text, if infix or true true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.a or .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "text, if infix or true false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.a or .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "text, if infix or false false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.a or .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": false, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "if with and true true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with and true false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "if with and false true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": false, "b": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "if with and false false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": false, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "if with or true true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with or true false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with or false true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": false, "b": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with or false false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": false, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "if with not true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if not .Values.a}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "if with not false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if not .Values.a}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with and 5 conditions all true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b .Values.c .Values.d .Values.e}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": true, "c": true, "d": true, "e": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with or 5 conditions all true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b .Values.c .Values.d .Values.e}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": true, "c": true, "d": true, "e": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with and infix",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.a and .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if else with or both true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 4, Indent: 2},
				{Type: types.TokenElse, Value: "{{else}}", Line: 5, Indent: 0},
				{Type: types.TokenText, Value: "  key3: val3\n", Line: 6, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 7, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key1: val1\n", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "  key2: val2\n", Line: 4, Indent: 2},
			},
		},
		{
			name: "if with and true false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if and .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
			},
		},
		{
			name: "if with or true false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if or .Values.a .Values.b}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": true, "b": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "if with not false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if not .Values.a}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{"a": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 3, Indent: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createCtx(tt.values)
			nodes, err := ast.ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			result, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluateASTEmbeddedIf(t *testing.T) {
	// Helper to create EvalContext with given values
	createCtx := func(values map[string]interface{}) *types.EvalContext {
		return eval.NewEvalContext(values, map[string]interface{}{"name": "test"})
	}

	tests := []struct {
		name     string
		tokens   []types.Token
		values   map[string]interface{}
		expected []types.Token
	}{
		{
			name: "outer true, inner true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 2},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 4},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 7, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
			values: map[string]interface{}{"outer": true, "inner": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 4},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
		},
		{
			name: "outer true, inner false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 2},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 4},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 7, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
			values: map[string]interface{}{"outer": true, "inner": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
		},
		{
			name: "outer false, inner true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 2},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 4},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 7, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
			values: map[string]interface{}{"outer": false, "inner": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
		},
		{
			name: "outer false, inner false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 2},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 4},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 7, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
			values: map[string]interface{}{"outer": false, "inner": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 8, Indent: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createCtx(tt.values)
			nodes, err := ast.ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			result, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}


func TestEvaluateAST_Range(t *testing.T) {
	// Helper to create EvalContext with given values
	createCtx := func(values map[string]interface{}) *types.EvalContext {
		return eval.NewEvalContext(values, map[string]interface{}{"name": "test"})
	}

	tests := []struct {
		name           string
		tokens         []types.Token
		values         map[string]interface{}
		expectedCount  int // expected number of output tokens
		checkContains  string // substring that should be in output
	}{
		{
			name: "range over simple list",
			tokens: []types.Token{
				{Type: types.TokenRange, Value: "{{range .Values.items}}", Line: 1, Indent: 0},
				{Type: types.TokenAction, Value: "{{.}}", Line: 2, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 3, Indent: 0},
			},
			values: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			expectedCount: 3,
			checkContains: "abc",
		},
		{
			name: "range over list of maps",
			tokens: []types.Token{
				{Type: types.TokenRange, Value: "{{range .Values.items}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "{{.name}}", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: ": ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "{{.value}}", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 2, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 3, Indent: 0},
			},
			values: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "item1", "value": "val1"},
					map[string]interface{}{"name": "item2", "value": "val2"},
				},
			},
			expectedCount: 10, // 5 tokens per iteration x 2 iterations
			checkContains: "item1: val1",
		},
		{
			name: "range over empty list",
			tokens: []types.Token{
				{Type: types.TokenRange, Value: "{{range .Values.items}}", Line: 1, Indent: 0},
				{Type: types.TokenAction, Value: "{{.}}", Line: 2, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 3, Indent: 0},
			},
			values: map[string]interface{}{
				"items": []interface{}{},
			},
			expectedCount: 0,
			checkContains: "",
		},
		{
			name: "range over map values",
			tokens: []types.Token{
				{Type: types.TokenRange, Value: "{{range .Values.config}}", Line: 1, Indent: 0},
				{Type: types.TokenAction, Value: "{{.}}", Line: 2, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 3, Indent: 0},
			},
			values: map[string]interface{}{
				"config": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			expectedCount: 2,
			checkContains: "value1",
		},
		{
			name: "range with if - only show enabled users",
			tokens: []types.Token{
				{Type: types.TokenRange, Value: "{{range .Values.users}}", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .enabled}}", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "  ", Line: 3, Indent: 4},
				{Type: types.TokenAction, Value: "{{.name}}", Line: 3, Indent: 4},
				{Type: types.TokenText, Value: ": enabled\n", Line: 3, Indent: 4},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "alice", "enabled": true},
					map[string]interface{}{"name": "bob", "enabled": false},
					map[string]interface{}{"name": "charlie", "enabled": true},
				},
			},
			expectedCount: 6, // 3 tokens per enabled user x 2 enabled users
			checkContains: "alice: enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createCtx(tt.values)
			nodes, err := ast.ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			result, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check expected token count
			if len(result) != tt.expectedCount {
				t.Errorf("expected %d tokens, got %d", tt.expectedCount, len(result))
			}

			// Check that output contains expected substring
			if tt.checkContains != "" {
				output := ""
				for _, tok := range result {
					output += tok.Value
				}
				if !contains(output, tt.checkContains) {
					t.Errorf("expected output to contain %q, got %q", tt.checkContains, output)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(substr) <= len(s) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestEvaluateASTElse(t *testing.T) {
	// Helper to create EvalContext with given values
	createCtx := func(values map[string]interface{}) *types.EvalContext {
		return eval.NewEvalContext(values, map[string]interface{}{"name": "test"})
	}

	tests := []struct {
		name     string
		tokens   []types.Token
		values   map[string]interface{}
		expected []types.Token
	}{
		{
			name: "if else end with static before and after, condition true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.cond}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "if content\n", Line: 3, Indent: 0},
				{Type: types.TokenElse, Value: "{{else}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "else content\n", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 7, Indent: 0},
			},
			values: map[string]interface{}{"cond": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "if content\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 7, Indent: 0},
			},
		},
		{
			name: "if else end with static before and after, condition false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.cond}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "if content\n", Line: 3, Indent: 0},
				{Type: types.TokenElse, Value: "{{else}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "else content\n", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 7, Indent: 0},
			},
			values: map[string]interface{}{"cond": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "else content\n", Line: 5, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 7, Indent: 0},
			},
		},
		{
			name: "embedded if if-else end with static before and after, outer true, inner true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "inner if content\n", Line: 5, Indent: 0},
				{Type: types.TokenElse, Value: "{{else}}", Line: 6, Indent: 0},
				{Type: types.TokenText, Value: "inner else content\n", Line: 7, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 8, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 9, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
			values: map[string]interface{}{"outer": true, "inner": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "inner if content\n", Line: 5, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
		},
		{
			name: "embedded if if-else end with static before and after, outer true, inner false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "inner if content\n", Line: 5, Indent: 0},
				{Type: types.TokenElse, Value: "{{else}}", Line: 6, Indent: 0},
				{Type: types.TokenText, Value: "inner else content\n", Line: 7, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 8, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 9, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
			values: map[string]interface{}{"outer": true, "inner": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "inner else content\n", Line: 7, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
		},
		{
			name: "embedded if if-end else-end with static before and after, outer true, inner true",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
				{Type: types.TokenElse, Value: "{{else}}", Line: 7, Indent: 0},
				{Type: types.TokenText, Value: "outer else content\n", Line: 8, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 9, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
			values: map[string]interface{}{"outer": true, "inner": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
		},
		{
			name: "embedded if if-end else-end with static before and after, outer false",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.outer}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "outer content\n", Line: 3, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.inner}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "inner content\n", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
				{Type: types.TokenElse, Value: "{{else}}", Line: 7, Indent: 0},
				{Type: types.TokenText, Value: "outer else content\n", Line: 8, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 9, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
			values: map[string]interface{}{"outer": false, "inner": true},
			expected: []types.Token{
				{Type: types.TokenText, Value: "static before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "outer else content\n", Line: 8, Indent: 0},
				{Type: types.TokenText, Value: "static after\n", Line: 10, Indent: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createCtx(tt.values)
			nodes, err := ast.ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			result, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluateAST_With(t *testing.T) {
	// Helper to create EvalContext with given values
	createCtx := func(values map[string]interface{}) *types.EvalContext {
		return eval.NewEvalContext(values, map[string]interface{}{"name": "test"})
	}

	tests := []struct {
		name     string
		tokens   []types.Token
		values   map[string]interface{}
		expected []types.Token
	}{
		{
			name: "with truthy value, rescopes context to nested field",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{.name}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{
				"config": map[string]interface{}{"name": "myconfig"},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "myconfig", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
			},
		},
		{
			name: "with missing value, body skipped",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  inside\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 5, Indent: 0},
			},
			values:   map[string]interface{}{},
			expected: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 5, Indent: 0},
			},
		},
		{
			name: "with nil value, body skipped",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  inside\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{"config": nil},
			expected: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 5, Indent: 0},
			},
		},
		{
			name: "with false value, body skipped",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  inside\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{"config": false},
			expected: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 5, Indent: 0},
			},
		},
		{
			name: "with else, truthy executes body",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  found: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{.name}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenElse, Value: "{{else}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "  not found\n", Line: 5, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
			},
			values: map[string]interface{}{
				"config": map[string]interface{}{"name": "myconfig"},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  found: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "myconfig", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
			},
		},
		{
			name: "with else, falsy executes else",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  found\n", Line: 3, Indent: 2},
				{Type: types.TokenElse, Value: "{{else}}", Line: 4, Indent: 0},
				{Type: types.TokenText, Value: "  not found\n", Line: 5, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
			},
			values:   map[string]interface{}{},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  not found\n", Line: 5, Indent: 2},
			},
		},
		{
			name: "with deep nested expression, rescopes correctly",
			tokens: []types.Token{
				{Type: types.TokenWith, Value: "{{with .Values.app.server}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  host: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "{{.host}}", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "\n  port: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "{{.port}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{
				"app": map[string]interface{}{
					"server": map[string]interface{}{
						"host": "localhost",
						"port": "8080",
					},
				},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "  host: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "localhost", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "\n  port: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "8080", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
			},
		},
		{
			name: "with nested inside if, both truthy",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.enabled}}", Line: 2, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 4, Indent: 2},
				{Type: types.TokenAction, Value: "{{.name}}", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 4, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
			},
			values: map[string]interface{}{
				"enabled": true,
				"config":  map[string]interface{}{"name": "myconfig"},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 4, Indent: 2},
				{Type: types.TokenAction, Value: "myconfig", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 4, Indent: 0},
			},
		},
		{
			name: "with nested inside if, if false skips with",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenIf, Value: "{{if .Values.enabled}}", Line: 2, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "  key: value\n", Line: 4, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 7, Indent: 0},
			},
			values: map[string]interface{}{
				"enabled": false,
				"config":  map[string]interface{}{"name": "myconfig"},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 7, Indent: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createCtx(tt.values)
			nodes, err := ast.ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			result, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluateAST_WithNested(t *testing.T) {
	createCtx := func(values map[string]interface{}) *types.EvalContext {
		return eval.NewEvalContext(values, map[string]interface{}{"name": "test"})
	}

	tests := []struct {
		name     string
		tokens   []types.Token
		values   map[string]interface{}
		expected []types.Token
	}{
		{
			name: "with inside with, both truthy",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.app}}", Line: 2, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .server}}", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "  host: ", Line: 4, Indent: 2},
				{Type: types.TokenAction, Value: "{{.host}}", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 4, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
			},
			values: map[string]interface{}{
				"app": map[string]interface{}{
					"server": map[string]interface{}{
						"host": "localhost",
					},
				},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  host: ", Line: 4, Indent: 2},
				{Type: types.TokenAction, Value: "localhost", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 4, Indent: 0},
			},
		},
		{
			name: "with inside with, inner falsy falls to else",
			tokens: []types.Token{
				{Type: types.TokenWith, Value: "{{with .Values.app}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  name: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "{{.name}}", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 2, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .server}}", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "  host: value\n", Line: 4, Indent: 2},
				{Type: types.TokenElse, Value: "{{else}}", Line: 5, Indent: 0},
				{Type: types.TokenText, Value: "  no server\n", Line: 6, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 7, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 8, Indent: 0},
			},
			values: map[string]interface{}{
				"app": map[string]interface{}{"name": "myapp"},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "  name: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "myapp", Line: 2, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  no server\n", Line: 6, Indent: 2},
			},
		},
		{
			name: "with inside with, outer falsy skips everything",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.app}}", Line: 2, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .server}}", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "  host: value\n", Line: 4, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 6, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 7, Indent: 0},
			},
			values:   map[string]interface{}{},
			expected: []types.Token{
				{Type: types.TokenText, Value: "before\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "after\n", Line: 7, Indent: 0},
			},
		},
		{
			name: "with inside range, rescopes per iteration",
			tokens: []types.Token{
				{Type: types.TokenRange, Value: "{{range .Values.items}}", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{.key}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"config": map[string]interface{}{"key": "val1"}},
					map[string]interface{}{"config": map[string]interface{}{"key": "val2"}},
				},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "val1", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "val2", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
			},
		},
		{
			name: "with inside range, some items missing key skips those",
			tokens: []types.Token{
				{Type: types.TokenRange, Value: "{{range .Values.items}}", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{.key}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"config": map[string]interface{}{"key": "val1"}},
					map[string]interface{}{"name": "noconfig"},
					map[string]interface{}{"config": map[string]interface{}{"key": "val3"}},
				},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "val1", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "  key: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "val3", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
			},
		},
		{
			name: "range inside with, with rescopes then range iterates",
			tokens: []types.Token{
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 1, Indent: 0},
				{Type: types.TokenRange, Value: "{{range .items}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "- ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{.}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{
				"config": map[string]interface{}{
					"items": []interface{}{"a", "b", "c"},
				},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "- ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "a", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "- ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "b", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
				{Type: types.TokenText, Value: "- ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "c", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 3, Indent: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createCtx(tt.values)
			nodes, err := ast.ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			result, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluateAST_RootContext(t *testing.T) {
	// Helper to create EvalContext with given values
	createCtx := func(values map[string]interface{}) *types.EvalContext {
		return eval.NewEvalContext(values, map[string]interface{}{"name": "test"})
	}

	tests := []struct {
		name     string
		tokens   []types.Token
		values   map[string]interface{}
		expected []types.Token
	}{
		{
			name: "using $ inside with to access root values",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  inside: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{.name}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n  from-root: ", Line: 4, Indent: 2},
				{Type: types.TokenAction, Value: "{{$.Values.global}}", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 4, Indent: 0},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 5, Indent: 0},
			},
			values: map[string]interface{}{
				"config": map[string]interface{}{"name": "myconfig"},
				"global": "root-value",
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "data:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  inside: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "myconfig", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "\n  from-root: ", Line: 4, Indent: 2},
				{Type: types.TokenAction, Value: "root-value", Line: 4, Indent: 2},
				{Type: types.TokenText, Value: "\n", Line: 4, Indent: 0},
			},
		},
		{
			name: "using $ inside range to access root values",
			tokens: []types.Token{
				{Type: types.TokenText, Value: "items:\n", Line: 1, Indent: 0},
				{Type: types.TokenRange, Value: "{{range .Values.items}}", Line: 2, Indent: 0},
				{Type: types.TokenText, Value: "  - ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{.}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: " (root: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "{{$.Values.prefix}}", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: ")\n", Line: 3, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 4, Indent: 0},
			},
			values: map[string]interface{}{
				"items":  []interface{}{"a", "b"},
				"prefix": "test",
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "items:\n", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "  - ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "a", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: " (root: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "test", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: ")\n", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: "  - ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "b", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: " (root: ", Line: 3, Indent: 2},
				{Type: types.TokenAction, Value: "test", Line: 3, Indent: 2},
				{Type: types.TokenText, Value: ")\n", Line: 3, Indent: 2},
			},
		},
		{
			name: "$ alone inside with returns the root map",
			tokens: []types.Token{
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 1, Indent: 0},
				{Type: types.TokenAction, Value: "{{$}}", Line: 2, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 3, Indent: 0},
			},
			values: map[string]interface{}{
				"config": map[string]interface{}{"name": "myconfig"},
				"extra":  "data",
			},
			expected: []types.Token{
				// The root map rendered as string representation
				// (the ActionNode uses fmt.Sprintf("%v", resultVal))
				{Type: types.TokenAction, Value: "map[config:map[name:myconfig] extra:data]", Line: 2, Indent: 2},
			},
		},
		{
			name: "$.Chart accessible from inside with",
			tokens: []types.Token{
				{Type: types.TokenWith, Value: "{{with .Values.config}}", Line: 1, Indent: 0},
				{Type: types.TokenText, Value: "chart: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "{{$.Chart.name}}", Line: 2, Indent: 2},
				{Type: types.TokenEnd, Value: "{{end}}", Line: 3, Indent: 0},
			},
			values: map[string]interface{}{
				"config": map[string]interface{}{"name": "myconfig"},
			},
			expected: []types.Token{
				{Type: types.TokenText, Value: "chart: ", Line: 2, Indent: 2},
				{Type: types.TokenAction, Value: "test", Line: 2, Indent: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createCtx(tt.values)
			nodes, err := ast.ParseAST(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error parsing AST: %v", err)
			}
			result, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}