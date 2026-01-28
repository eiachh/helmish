package helmishlib

import (
	"helmish/internal/renderer"
)

// Chart represents the Helm chart options
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

// Render calls the internal renderer to render the chart
func Render(opts Options) (map[string]string, error) {
	loadedProfile, err := loadProfile(opts.Profile.Name)
	if err != nil {
		return nil, err
	}
	internalOpts := renderer.Options{
		Chart: renderer.Chart{
			Path: opts.Chart.Path,
		},
		Profile: loadedProfile,
	}
	return renderer.RenderChart(internalOpts)
}