package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"helmish/pkg/helmishlib"
)

func renderToString(tokensList [][]helmishlib.Token) string {
	var parts []string
	for _, tokens := range tokensList {
		lineMap := make(map[int][]string)
		lineIndents := make(map[int]int)
		for _, token := range tokens {
			if _, ok := lineIndents[token.Line]; !ok {
				lineIndents[token.Line] = token.Indent
			}
			lineMap[token.Line] = append(lineMap[token.Line], token.Type.String()+": "+token.Value)
		}
		var sortedLines []int
		for line := range lineMap {
			sortedLines = append(sortedLines, line)
		}
		sort.Ints(sortedLines)
		var lines []string
		for _, line := range sortedLines {
			indent := lineIndents[line]
			parts := lineMap[line]
			if len(parts) > 0 {
				parts[0] = strings.Repeat(" ", indent) + parts[0]
			}
			lines = append(lines, strings.Join(parts, ""))
		}
		parts = append(parts, strings.Join(lines, "\n"))
	}
	return strings.Join(parts, "\n---\n")
}

func renderRawToString(tokensList [][]helmishlib.Token) string {
	var parts []string
	for _, tokens := range tokensList {
		var raw []string
		for _, token := range tokens {
			raw = append(raw, token.Value)
		}
		parts = append(parts, strings.Join(raw, ""))
	}
	return strings.Join(parts, "\n---\n")
}

type model struct {
	rawContent map[string]string
	rendered   map[string][][]helmishlib.Token
	showPopup  bool
	selected   int
	files      []string
	width      int
	height     int
	currentFile string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 3
	case tea.KeyMsg:
		if m.showPopup {
			switch msg.String() {
			case "up":
				if m.selected > 0 {
					m.selected--
				}
			case "down":
				if m.selected < len(m.files)-1 {
					m.selected++
				}
			case "enter":
				if m.selected >= 0 && m.selected < len(m.files) {
					m.currentFile = m.files[m.selected]
				}
				m.showPopup = false
			case "esc":
				m.showPopup = false
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		} else {
			switch msg.String() {
			case "0":
				m.showPopup = true
				m.selected = 0
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.showPopup {
		return m.renderOverlay()
	}
	return m.renderMain()
}

func (m model) renderMain() string {
	halfWidth := m.width / 2
	// Adjust for border (2) and padding (2) on each side
	adjustedWidth := halfWidth - 4
	leftStyle := lipgloss.NewStyle().
		Width(adjustedWidth).
		Height(m.height).
		Border(lipgloss.ThickBorder()).
		Padding(1)

	rightStyle := lipgloss.NewStyle().
		Width(adjustedWidth).
		Height(m.height).
		Border(lipgloss.ThickBorder()).
		Padding(1)

	// Left panel: raw content
	leftContent := m.rawContent[m.currentFile]
	left := leftStyle.Render(leftContent)

	// Right panel: rendered content
	rightContent := renderToString(m.rendered[m.currentFile])
	right := rightStyle.Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m model) renderOverlay() string {
	background := lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, m.renderMain())
	popup := m.renderFileSelectorPopup()
	// Popup dimensions: width = 50 + 2(border) + 2(padding) = 54
	// height = 20 + 2(border) + 2(padding) = 24
	fgWidth := 54
	fgHeight := 24
	x := (m.width - fgWidth) / 2
	y := (m.height - fgHeight) / 2
	// Ensure it fits within the screen
	if x < 0 {
		x = 0
	}
	if x+fgWidth > m.width {
		x = m.width - fgWidth
	}
	if y < 0 {
		y = 0
	}
	if y+fgHeight > m.height {
		y = m.height - fgHeight
	}
	return overlayStrings(background, popup, m.width, m.height, fgWidth, fgHeight, x, y)
}

func overlayStrings(bg, fg string, width, height, fgWidth, fgHeight, x, y int) string {
	bgLines := strings.Split(bg, "\n")
	fgLines := strings.Split(fg, "\n")
	for i := 0; i < len(bgLines) && i < height; i++ {
		if i >= y && i < y+fgHeight && (i-y) < len(fgLines) {
			fgLine := fgLines[i-y]
			if len(fgLine) > 0 {
				start := x
				if start < 0 {
					start = 0
				}
				end := start + len(fgLine)
				if end > len(bgLines[i]) {
					end = len(bgLines[i])
				}
				if start < len(bgLines[i]) && end > start {
					bgLines[i] = bgLines[i][:start] + fgLine[:end-start] + bgLines[i][end:]
				}
			}
		}
	}
	return strings.Join(bgLines, "\n")
}

func (m model) renderFileSelectorPopup() string {
	popupStyle := lipgloss.NewStyle().
		Width(50).
		Height(20).
		Border(lipgloss.ThickBorder()).
		Padding(1)

	content := "File Selector (use arrow keys, enter to select, esc to cancel):\n\n"

	// Find the current displayed file index
	currentIndex := -1
	for i, f := range m.files {
		if f == m.currentFile {
			currentIndex = i
			break
		}
	}

	for i, f := range m.files {
		prefix := "  "
		if i == m.selected {
			prefix = "> "
		}
		if i == currentIndex {
			prefix = "* "
		}
		if i == m.selected && i == currentIndex {
			prefix = "*>"
		}
		content += prefix + f + "\n"
	}

	return popupStyle.Render(content)
}

func main() {
	opts := parseConfig()

	// Load the chart
	helmish, err := helmishlib.NewHelmish(opts.Chart.Path)
	if err != nil {
		log.Fatalf("Error loading chart: %v", err)
	}

	// Render the chart
	rendered, err := helmish.Render(opts.Profile)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare raw content
	rawContent := make(map[string]string)
	for filename, templates := range rendered {
		rawContent[filename] = renderRawToString(templates)
	}

	// Print rendered
	for filename, tokensList := range rendered {
		fmt.Printf("File: %s\n", filename)
		fmt.Println(renderToString(tokensList))
		fmt.Println("---")
	}
}