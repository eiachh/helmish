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