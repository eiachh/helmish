package helmishlib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenIfExamples(t *testing.T) {
	examplesRoot := filepath.Join("..", "..", "examples", "if-examples")
	testdataRoot := filepath.Join("..", "..", "testdata", "examples", "if-examples")

	entries, err := os.ReadDir(examplesRoot)
	if err != nil {
		t.Fatalf("failed to read examples: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		chartPath := filepath.Join(examplesRoot, e.Name())
		chartYAML := filepath.Join(chartPath, "Chart.yaml")
		if _, err := os.Stat(chartYAML); err != nil {
			continue
		}

		t.Run(e.Name(), func(t *testing.T) {
			h, err := NewHelmish(chartPath)
			if err != nil {
				t.Fatalf("NewHelmish: %v", err)
			}
			tokens, err := h.Render(Profile{Name: "default"})
			if err != nil {
				t.Fatalf("Render: %v", err)
			}
			rendered := RenderAllFilesToString(tokens)

			for filename, got := range rendered {
				goldenPath := filepath.Join(testdataRoot, e.Name(), filename)
				wantBytes, err := os.ReadFile(goldenPath)
				if err != nil {
					t.Skipf("no golden for %s: %v", filename, err)
					continue
				}
				want := string(wantBytes)
				if got != want {
					t.Errorf("mismatch %s:\n--- got ---\n%s\n--- want ---\n%s", filename, got, want)
				}
			}
		})
	}
}

func TestGoldenEdgeCases(t *testing.T) {
	examplesRoot := filepath.Join("..", "..", "examples", "edge-cases")
	testdataRoot := filepath.Join("..", "..", "testdata", "examples", "edge-cases")

	entries, err := os.ReadDir(examplesRoot)
	if err != nil {
		t.Fatalf("failed to read examples: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		chartPath := filepath.Join(examplesRoot, e.Name())
		chartYAML := filepath.Join(chartPath, "Chart.yaml")
		if _, err := os.Stat(chartYAML); err != nil {
			continue
		}

		t.Run(e.Name(), func(t *testing.T) {
			h, err := NewHelmish(chartPath)
			if err != nil {
				t.Fatalf("NewHelmish: %v", err)
			}
			tokens, err := h.Render(Profile{Name: "default"})
			if err != nil {
				t.Fatalf("Render: %v", err)
			}
			rendered := RenderAllFilesToString(tokens)

			for filename, got := range rendered {
				goldenPath := filepath.Join(testdataRoot, e.Name(), filename)
				wantBytes, err := os.ReadFile(goldenPath)
				if err != nil {
					t.Skipf("no golden for %s: %v", filename, err)
					continue
				}
				want := string(wantBytes)
				if got != want {
					t.Errorf("mismatch %s:\n--- got ---\n%s\n--- want ---\n%s", filename, got, want)
				}
			}
		})
	}
}

func TestHeaderFormat(t *testing.T) {
	// Test that headers are correctly formatted with --- and # Source:
	chartPath := filepath.Join("..", "..", "examples", "edge-cases", "multi-template-chart")
	h, err := NewHelmish(chartPath)
	if err != nil {
		t.Fatalf("NewHelmish: %v", err)
	}
	tokens, err := h.Render(Profile{Name: "default"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	// Check that each file has the correct header structure
	for filename, fileTokens := range tokens {
		if len(fileTokens) < 2 {
			t.Errorf("file %s: expected at least 2 token lines (header), got %d", filename, len(fileTokens))
			continue
		}

		// First line should be "---"
		firstLine := fileTokens[0]
		if len(firstLine) != 1 || firstLine[0].Value != "---" {
			t.Errorf("file %s: first line should be '---', got %v", filename, firstLine)
		}

		// Second line should be "# Source: <chart>/templates/<filename>"
		secondLine := fileTokens[1]
		if len(secondLine) != 1 {
			t.Errorf("file %s: second line should have 1 token, got %d", filename, len(secondLine))
			continue
		}
		expectedSource := "# Source: multi-template-chart/templates/" + filename
		if secondLine[0].Value != expectedSource {
			t.Errorf("file %s: second line should be '%s', got '%s'", filename, expectedSource, secondLine[0].Value)
		}
	}
}

func TestMissingChartName(t *testing.T) {
	// Test that charts without a Name field fall back to "unknown"
	chartPath := filepath.Join("..", "..", "examples", "edge-cases", "missing-name")
	h, err := NewHelmish(chartPath)
	if err != nil {
		t.Fatalf("NewHelmish: %v", err)
	}
	tokens, err := h.Render(Profile{Name: "default"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	fileTokens, ok := tokens["configmap.yaml"]
	if !ok {
		t.Fatal("expected configmap.yaml in output")
	}

	if len(fileTokens) < 2 {
		t.Fatalf("expected at least 2 token lines, got %d", len(fileTokens))
	}

	// Second line should contain "unknown" as the chart name
	secondLine := fileTokens[1]
	if len(secondLine) != 1 {
		t.Fatalf("second line should have 1 token, got %d", len(secondLine))
	}

	expectedSource := "# Source: unknown/templates/configmap.yaml"
	if secondLine[0].Value != expectedSource {
		t.Errorf("expected '%s', got '%s'", expectedSource, secondLine[0].Value)
	}
}
