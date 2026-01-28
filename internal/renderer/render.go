package renderer

import (
	"os"
	"path/filepath"
)

// Values holds the chart values files (filename -> content)
type Values map[string]string

// Metadata holds the chart metadata files (filename -> content)
type Metadata map[string]string

// YamlTemplates holds the YAML template files (filename -> content)
type YamlTemplates map[string]string

// TplFiles holds the template files (filename -> content)
type TplFiles map[string]string

// Chart represents the Helm chart data
type Chart struct {
	Path          string
	Values        Values
	Metadata      Metadata
	YamlTemplates YamlTemplates
	TplFiles      TplFiles
}

// Profile represents the profile options
type Profile struct {
	Name          string
	Capabilities  Capabilities
	// Add more fields as needed
}

// Capabilities represents Helm capabilities
type Capabilities struct {
	KubeVersion string
	APIVersions []string
}

// Options holds the options for rendering
type Options struct {
	Chart   Chart
	Profile Profile
}

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
func RenderChart(opts Options) (map[string]string, error) {
	// For now, return the first element of YamlTemplates
	var firstContent string
	var firstFilename string
	for filename, content := range opts.Chart.YamlTemplates {
		firstContent = content
		firstFilename = filename
		break
	}
	if firstContent == "" {
		// If no templates, return empty
		return map[string]string{}, nil
	}
	return map[string]string{firstFilename: firstContent}, nil
}