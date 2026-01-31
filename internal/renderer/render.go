package renderer

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"helmish/internal/renderer/ast"
	"helmish/internal/renderer/eval"
	"helmish/internal/renderer/parser"
	"helmish/internal/renderer/tokenizer"
	"helmish/internal/renderer/types"
)

// Aliases for public API
type ValueData = types.ValueData
type Chart = types.Chart
type Profile = types.Profile
type Capabilities = types.Capabilities
type Options = types.Options

// LoadChart loads the chart from the given path
func LoadChart(path string) (types.Chart, error) {
	chart := types.Chart{
		Path:          path,
		Values:        make(types.Values),
		Metadata:      make(types.Metadata),
		YamlTemplates: make(types.YamlTemplates),
		TplFiles:      make(types.TplFiles),
	}

	// Load Chart.yaml for metadata
	chartYamlPath := filepath.Join(path, "Chart.yaml")
	if content, err := os.ReadFile(chartYamlPath); err == nil {
		var parsed interface{}
		if err := yaml.Unmarshal(content, &parsed); err != nil {
			return chart, err
		}
		chart.Metadata["Chart.yaml"] = types.ValueData{
			Raw:    string(content),
			Parsed: parsed,
		}
	} else {
		return chart, err
	}

	// Load values.yaml for values
	valuesYamlPath := filepath.Join(path, "values.yaml")
	if content, err := os.ReadFile(valuesYamlPath); err == nil {
		var parsed interface{}
		if err := yaml.Unmarshal(content, &parsed); err != nil {
			return chart, err
		}
		chart.Values["values.yaml"] = types.ValueData{
			Raw:    string(content),
			Parsed: parsed,
		}
	} else {
		return chart, err
	}

	// Load templates
	templatesPath := filepath.Join(path, "templates")
	err := filepath.Walk(templatesPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(templatesPath, p)
		content, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		ext := filepath.Ext(p)
		switch ext {
		case ".yaml", ".yml":
			chart.YamlTemplates[relPath] = string(content)
		case ".tpl":
			chart.TplFiles[relPath] = string(content)
		}
		return nil
	})
	if err != nil {
		// Return error if templates directory doesn't exist or any read error
		return chart, err
	}

	return chart, nil
}

// RenderChart renders the Helm chart using the TUI
func RenderChart(opts Options) (map[string][][]types.Token, error) {
	result := make(map[string][][]types.Token)
	// Create eval context from values and chart
	var values, chart interface{}
	if val, ok := opts.Chart.Values["values.yaml"]; ok {
		values = val.Parsed
	}
	if meta, ok := opts.Chart.Metadata["Chart.yaml"]; ok {
		chart = meta.Parsed
	}
	ctx := eval.NewEvalContext(values, chart)
	for filename, content := range opts.Chart.YamlTemplates {
		blocks := parser.CollectBlocks(content)
		for _, doc := range blocks {
			tokens := tokens.Tokenize(doc)
			// Parse AST from tokens
			nodes, err := ast.ParseAST(tokens)
			if err != nil {
				return nil, err
			}
			// Evaluate the AST
			evaluatedTokens, err := eval.EvaluateAST(nodes, ctx)
			if err != nil {
				return nil, err
			}
			result[filename] = append(result[filename], evaluatedTokens)
		}
	}
	return result, nil
}