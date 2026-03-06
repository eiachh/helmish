package helmishlib

import (
	"helmish/internal/renderer"
	"helmish/internal/renderer/types"
)

// ValueData holds the raw content and parsed data for a values file
type ValueData = renderer.ValueData

// KeyValueBlock represents a YAML key-value pair
type KeyValueBlock = types.KeyValueBlock

// TemplateBlock represents a Helm template block
type TemplateBlock = types.TemplateBlock

// Block represents a single line block in a rendered template
type Block = types.Block

// TokenType represents the type of token
type TokenType = types.TokenType

// Token represents a single token in the template
type Token = types.Token

// TokenText is a constant for text tokens
const TokenText = types.TokenText

// TokenIf is a constant for if tokens
const TokenIf = types.TokenIf

// TokenElse is a constant for else tokens
const TokenElse = types.TokenElse

// TokenEnd is a constant for end tokens
const TokenEnd = types.TokenEnd

// TokenRange is a constant for range tokens
const TokenRange = types.TokenRange

// TokenAction is a constant for action tokens
const TokenAction = types.TokenAction

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
func (h *Helmish) Render(profile Profile) (map[string][][]Token, error) {
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

// RenderTokensToString converts a 2D slice of tokens to a string representation.
// Each inner slice represents a line (or document), and tokens are concatenated
// to form the rendered output. Newlines are added between lines.
func RenderTokensToString(tokens [][]Token) string {
	var result string
	for i, line := range tokens {
		if i > 0 {
			result += "\n"
		}
		for _, tok := range line {
			result += tok.Value
		}
	}
	return result
}

// RenderAllFilesToString converts a map of files (each containing 2D token slices)
// to a map of rendered strings.
func RenderAllFilesToString(fileTokens map[string][][]Token) map[string]string {
	result := make(map[string]string)
	for filename, tokens := range fileTokens {
		result[filename] = RenderTokensToString(tokens)
	}
	return result
}