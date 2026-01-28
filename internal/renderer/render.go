package renderer

// Chart represents the Helm chart options
type Chart struct {
	Path string
	// Add more fields as needed
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

// RenderChart renders the Helm chart using the TUI
func RenderChart(opts Options) (map[string]string, error) {
	// Hardcoded rendered content for now
	// Use opts.Profile in future rendering
	rendered := map[string]string{
		"template.yaml": "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: example\n",
		"values.yaml":   "key: value\n",
	}
	return rendered, nil
}