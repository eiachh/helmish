package renderer

import (
	"os"
	"path/filepath"
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
		chart.Metadata["Chart.yaml"] = string(content)
	} else {
		return chart, err // or handle error, but for now return error if can't read
	}

	// Load values.yaml for values
	valuesYamlPath := filepath.Join(path, "values.yaml")
	if content, err := os.ReadFile(valuesYamlPath); err == nil {
		chart.Values["values.yaml"] = string(content)
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
func RenderChart(opts Options) (map[string][]RenderedTemplate, error) {
	result := make(map[string][]RenderedTemplate)
	for filename, content := range opts.Chart.YamlTemplates {
		rendered := parseContent(content)
		result[filename] = rendered
	}
	return result, nil
}