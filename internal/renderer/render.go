package renderer

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadChart loads the chart from the given path
func LoadChart(path string) (Chart, error) {
	chart := Chart{
		Path:          path,
		Values:        make(Values),
		Metadata:      make(Metadata),
		YamlTemplates: make(YamlTemplates),
		TplFiles:      make(TplFiles),
	}

	// Load Chart.yaml for metadata
	chartYamlPath := filepath.Join(path, "Chart.yaml")
	if content, err := os.ReadFile(chartYamlPath); err == nil {
		var parsed interface{}
		if err := yaml.Unmarshal(content, &parsed); err != nil {
			return chart, err
		}
		chart.Metadata["Chart.yaml"] = ValueData{
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
		chart.Values["values.yaml"] = ValueData{
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
func RenderChart(opts Options) (map[string][][]Token, error) {
	result := make(map[string][][]Token)
	// Create eval context from values and chart
	var values, chart interface{}
	if val, ok := opts.Chart.Values["values.yaml"]; ok {
		values = val.Parsed
	}
	if meta, ok := opts.Chart.Metadata["Chart.yaml"]; ok {
		chart = meta.Parsed
	}
	ctx := NewEvalContext(values, chart)
	for filename, content := range opts.Chart.YamlTemplates {
		rendered := parseContent(content)
		for _, rt := range rendered {
			tokens := Tokenize(rt)
			// Evaluate the tokens
			evaluatedTokens, err := EvaluateTokens(tokens, ctx)
			if err != nil {
				return nil, err
			}
			result[filename] = append(result[filename], evaluatedTokens)
		}
	}
	return result, nil
}