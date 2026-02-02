package parser

import (
	"strings"

	"helmish/internal/renderer/types"
)

// collectMultilineTemplate collects lines from start until }} is found
func collectMultilineTemplate(lines []string, start int) (string, int) {
	var templateLines []string
	templateLines = append(templateLines, lines[start])
	i := start + 1
	for i < len(lines) && !strings.Contains(lines[i], "}}") {
		templateLines = append(templateLines, lines[i])
		i++
	}
	if i < len(lines) {
		templateLines = append(templateLines, lines[i])
		i++
	}
	content := strings.Join(templateLines, "\n")
	return content, i
}

// collectBlocks parses the content of a YAML file into a list of DocumentBlocks
func CollectBlocks(content string) []types.DocumentBlocks {
	lines := strings.Split(content, "\n")
	var blocks []types.DocumentBlocks
	var current types.DocumentBlocks
	i := 0
	for i < len(lines) {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			if len(current.Blocks) > 0 {
				blocks = append(blocks, current)
				current = types.DocumentBlocks{}
			}
			i++
			continue
		}
		if strings.Contains(line, "{{") && !strings.Contains(line, "}}") {
			// Multiline template block
			tmplContent, newI := collectMultilineTemplate(lines, i)
			indent := len(line) - len(strings.TrimLeft(line, " "))
			block := types.Block{
				Line:     i + 1,
				Type:    types.TemplateBlockType,
				Content: &types.TemplateBlock{RawContent: tmplContent},
				Indent:  indent,
			}
			current.Blocks = append(current.Blocks, block)
			i = newI
		} else {
			// Single line block
			indent := len(line) - len(strings.TrimLeft(line, " "))
			var block types.Block
			block.Line = i + 1
			block.Indent = indent
			if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
				block.Type = types.TemplateBlockType
				block.Content = &types.TemplateBlock{RawContent: line}
			} else {
				colonIndex := strings.Index(line, ":")
				if colonIndex != -1 {
					key := strings.TrimSpace(line[:colonIndex])
					valuePart := line[colonIndex+1:]
					value := strings.TrimSpace(valuePart)
					block.Type = types.KeyValueBlockType
					block.Content = &types.KeyValueBlock{Key: key, Value: value}
				} else {
					key := strings.TrimSpace(line)
					block.Type = types.KeyValueBlockType
					block.Content = &types.KeyValueBlock{Key: key, Value: ""}
				}
			}
			current.Blocks = append(current.Blocks, block)
			i++
		}
	}
	if len(current.Blocks) > 0 {
		blocks = append(blocks, current)
	}
	return blocks
}