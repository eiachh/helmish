package helmishlib

import (
	"helmish/internal/renderer"
)

// BlockContent represents content that can be raw or rendered
type BlockContent = renderer.BlockContent

// YamlKey represents a YAML key (possibly indented)
type YamlKey = renderer.YamlKey

// YamlKeyValue represents a YAML key-value pair
type YamlKeyValue = renderer.YamlKeyValue

// TemplateBlock represents a Helm template block
type TemplateBlock = renderer.TemplateBlock

// Block represents a single line block in a rendered template
type Block = renderer.Block

// RenderedTemplate represents a single YAML document with its blocks
type RenderedTemplate = renderer.RenderedTemplate

// Chart represents the Helm chart data
type Chart = renderer.Chart

// Profile represents the profile options (public, minimal)
type Profile struct {
	Name string
}

// Options holds the options for rendering (public)
type Options struct {
	Chart   Chart
	Profile Profile
}

// Helmish is the main library struct that holds the loaded chart
type Helmish struct {
	chart renderer.Chart
}

// NewHelmish creates a new Helmish instance by loading the chart from the given path
func NewHelmish(chartPath string) (*Helmish, error) {
	chart, err := renderer.LoadChart(chartPath)
	if err != nil {
		return nil, err
	}
	return &Helmish{
		chart: chart,
	}, nil
}

// loadProfile loads the profile from the default folder based on name
func loadProfile(name string) (renderer.Profile, error) {
	// Skeleton: hardcoded for now
	// In the future, look in default folder for profile file
	profile := renderer.Profile{
		Name: name,
		Capabilities: renderer.Capabilities{
			KubeVersion: "1.25",
			APIVersions: []string{"v1", "apps/v1"},
		},
	}
	return profile, nil
}

// Render calls the internal renderer to render the chart using the loaded chart
func (h *Helmish) Render(profile Profile) (map[string][]RenderedTemplate, error) {
	loadedProfile, err := loadProfile(profile.Name)
	if err != nil {
		return nil, err
	}
	internalOpts := renderer.Options{
		Chart:   h.chart,
		Profile: loadedProfile,
	}
	return renderer.RenderChart(internalOpts)
}