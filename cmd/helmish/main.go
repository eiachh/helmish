package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"helmish/pkg/helmishlib"
)

type model struct {
	left     string
	right    string
	rendered map[string]string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	leftStyle := lipgloss.NewStyle().
		Width(40).
		Height(10).
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	rightStyle := lipgloss.NewStyle().
		Width(40).
		Height(10).
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	// Left panel: small Helm template
	leftContent := "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: {{ .Chart.Name }}\ndata:\n  key: {{ .Values.key }}"
	left := leftStyle.Render("Helm Template\n" + leftContent)

	// Right panel: rendered content (show only the first file for now)
	var rightContent string
	for filename, content := range m.rendered {
		rightContent = filename + ":\n" + content
		break // Show only the first one
	}
	right := rightStyle.Render("Rendered Output\n" + rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}


func main() {
	// Parse config to get options
	opts := parseConfig()

	// Call the render function
	rendered, err := helmishlib.Render(opts)
	if err != nil {
		log.Fatal(err)
	}

	// Start the TUI
	m := model{
		left:     "",
		right:    "",
		rendered: rendered,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}