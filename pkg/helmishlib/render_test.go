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
