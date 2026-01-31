package renderer

import "strings"

// BlockContent represents content that can be raw or rendered
type BlockContent interface {
	Raw() string
	Rendered() string
}

// YamlKeyValue represents a YAML key-value pair (or key only if value is empty)
type YamlKeyValue struct {
	Key   string
	Value string
}

// Raw returns the raw YAML key-value pair
func (y YamlKeyValue) Raw() string {
	if y.Value == "" {
		return y.Key
	}
	return y.Key + ": " + y.Value
}

// Rendered returns the rendered YAML key-value pair
func (y YamlKeyValue) Rendered() string {
	if y.Value == "" {
		return y.Key
	}
	if strings.Contains(y.Value, "{{") && strings.Contains(y.Value, "}}") {
		// Simulate rendering by wrapping in [rendered ...]
		return y.Key + ": [rendered " + y.Value + "]"
	}
	return y.Raw()
}

// TemplateBlock represents a Helm template block
type TemplateBlock struct {
	RawContent string
}

// Raw returns the raw template content
func (t TemplateBlock) Raw() string {
	return t.RawContent
}

// Rendered returns the rendered template content
func (t TemplateBlock) Rendered() string {
	// Simulate rendering by replacing {{ ... }} with [rendered ...]
	result := strings.ReplaceAll(t.RawContent, "{{", "[rendered ")
	result = strings.ReplaceAll(result, "}}", "]")
	return result
}

// BlockType represents the type of a block
type BlockType int

const (
	YamlKeyValueBlock BlockType = iota
	TemplateBlockType
)

// Block represents a single line block in a rendered template
type Block struct {
	Line    int
	Type   BlockType
	Content BlockContent
	Indent int
}

// Raw returns the raw content of the block
func (b Block) Raw() string {
	return b.Content.Raw()
}

// Rendered returns the rendered content of the block
func (b Block) Rendered() string {
	return b.Content.Rendered()
}

// GetYamlKeyValue returns the YamlKeyValue if the block is of that type
func (b Block) GetYamlKeyValue() (*YamlKeyValue, bool) {
	if b.Type == YamlKeyValueBlock {
		if ykv, ok := b.Content.(*YamlKeyValue); ok {
			return ykv, true
		}
	}
	return nil, false
}

// GetTemplate returns the TemplateBlock if the block is of that type
func (b Block) GetTemplate() (*TemplateBlock, bool) {
	if b.Type == TemplateBlockType {
		if tb, ok := b.Content.(*TemplateBlock); ok {
			return tb, true
		}
	}
	return nil, false
}

// RenderedTemplate represents a single YAML document with its blocks
type RenderedTemplate struct {
	Blocks []Block
}

// ValueData holds the raw content and parsed data for a values file
type ValueData struct {
	Raw    string
	Parsed interface{}
}

// Values holds the chart values files (filename -> ValueData)
type Values map[string]ValueData

// Metadata holds the chart metadata files (filename -> ValueData)
type Metadata map[string]ValueData

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