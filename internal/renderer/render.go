package renderer

import (
	"os"
	"path/filepath"
	"strings"

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
			Parsed: normalizeChartMetadata(parsed),
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

	// Extract chart name for header
	chartName := getChartName(opts.Chart)

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
			// Prepend header tokens: --- and # Source: <chart>/templates/<filename>
			// Each header line is a separate []Token slice (each slice = one line)
			result[filename] = append(result[filename],
				[]types.Token{{Type: types.TokenText, Value: "---"}},
				[]types.Token{{Type: types.TokenText, Value: "# Source: " + chartName + "/templates/" + filename}},
				evaluatedTokens,
			)
		}
	}
	return result, nil
}

// getChartName extracts the chart name from metadata
func getChartName(chart types.Chart) string {
	if meta, ok := chart.Metadata["Chart.yaml"]; ok {
		if m, ok := meta.Parsed.(map[string]interface{}); ok {
			if name, ok := m["Name"].(string); ok {
				return name
			}
		}
		if m, ok := meta.Parsed.(map[interface{}]interface{}); ok {
			if name, ok := m["Name"].(string); ok {
				return name
			}
		}
	}
	return "unknown"
}

// normalizeChartMetadata converts map keys from YAML (lowercase) to Helm casing (Title)
func normalizeChartMetadata(v interface{}) interface{} {
	switch m := v.(type) {
	case map[interface{}]interface{}:
		out := make(map[interface{}]interface{}, len(m))
		for k, val := range m {
			key := normalizeKey(k)
			out[key] = normalizeChartMetadata(val)
		}
		return out
	case map[string]interface{}:
		out := make(map[string]interface{}, len(m))
		for k, val := range m {
			key := normalizeKey(k)
			out[key.(string)] = normalizeChartMetadata(val)
		}
		return out
	case []interface{}:
		for i, e := range m {
			m[i] = normalizeChartMetadata(e)
		}
		return m
	default:
		return v
	}
}

func normalizeKey(k interface{}) interface{} {
	s, ok := k.(string)
	if !ok {
		return k
	}
	switch strings.ToLower(s) {
	case "name":
		return "Name"
	case "home":
		return "Home"
	case "sources":
		return "Sources"
	case "version":
		return "Version"
	case "description":
		return "Description"
	case "keywords":
		return "Keywords"
	case "maintainers":
		return "Maintainers"
	case "engine":
		return "Engine"
	case "icon":
		return "Icon"
	case "appversion":
		return "AppVersion"
	case "deprecated":
		return "Deprecated"
	case "annotations":
		return "Annotations"
	case "kubeversion":
		return "KubeVersion"
	case "dependencies":
		return "Dependencies"
	case "type":
		return "Type"
	default:
		if len(s) > 0 {
			return strings.ToUpper(s[:1]) + s[1:]
		}
		return s
	}
}