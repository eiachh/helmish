package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"helmish/pkg/helmishlib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: helmish <chart-path>")
		os.Exit(1)
	}

	chartPath := os.Args[1]

	// Check if the path is absolute, if not make it relative to current directory
	if !filepath.IsAbs(chartPath) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		chartPath = filepath.Join(cwd, chartPath)
	}

	// Load the chart
	h, err := helmishlib.NewHelmish(chartPath)
	if err != nil {
		fmt.Printf("Error loading chart: %v\n", err)
		os.Exit(1)
	}

	// Render the chart
	tokens, err := h.Render(helmishlib.Profile{Name: "default"})
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
		os.Exit(1)
	}

	// Display tokenized output - group tokens by line
	fmt.Println("=== TOKENIZED OUTPUT ===")
	for filename, fileTokens := range tokens {
		fmt.Printf("\n--- %s ---\n", filename)
		for _, docTokens := range fileTokens {
			printTokensByLine(docTokens)
		}
	}

	// Display string rendered output
	fmt.Println("\n=== RENDERED OUTPUT ===")
	renderedFiles := helmishlib.RenderAllFilesToString(tokens)
	for filename, content := range renderedFiles {
		fmt.Printf("\n--- %s ---\n", filename)
		fmt.Println(strings.TrimSuffix(content, "\n"))
	}
}

// printTokensByLine prints tokens grouped by their line number
func printTokensByLine(tokens []helmishlib.Token) {
	if len(tokens) == 0 {
		return
	}

	// Group tokens by line
	lineGroups := make(map[int][]helmishlib.Token)
	var lineOrder []int
	
	for _, tok := range tokens {
		line := tok.Line
		if _, exists := lineGroups[line]; !exists {
			lineOrder = append(lineOrder, line)
		}
		lineGroups[line] = append(lineGroups[line], tok)
	}

	// Print tokens for each line
	for _, line := range lineOrder {
		group := lineGroups[line]
		// Filter out newline-only text tokens for display
		var filtered []helmishlib.Token
		for _, tok := range group {
			// Skip text tokens that are only newlines or empty
			if tok.Type == helmishlib.TokenText {
				val := strings.TrimSuffix(tok.Value, "\n")
				if val == "" {
					continue
				}
			}
			filtered = append(filtered, tok)
		}
		
		// Skip lines with no visible tokens
		if len(filtered) == 0 {
			continue
		}
		
		fmt.Printf("    Line %d: ", line)
		for i, tok := range filtered {
			// Trim newlines from displayed value for cleaner output
			val := strings.TrimSuffix(tok.Value, "\n")
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%s(%q)", tok.Type.String(), val)
		}
		fmt.Println()
	}
}
