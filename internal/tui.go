package internal

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	messages []string
	input    string
	loading  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if strings.TrimSpace(m.input) != "" {
				question := m.input
				m.messages = append(m.messages, question)
				m.input = ""
				m.loading = true
				return m, requestChat(question)
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	case responseMsg:
		m.messages = append(m.messages, msg.text)
		m.loading = false

	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500"))
	b.WriteString(headerStyle.Render("Codemate - Type & Press Enter\n"))
	b.WriteString("\n")

	for _, msg := range m.messages {
		b.WriteString(msg + "\n")
	}
	b.WriteString("\n")

	if m.loading {
		b.WriteString("Thinking...\n")
	}

	inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	b.WriteString(inputStyle.Render("$ " + m.input))

	return b.String()
}

type responseMsg struct{ text string }

func requestChat(prompt string) tea.Cmd {
	return func() tea.Msg {
		context, _ := GetProjectContext()
		return responseMsg{text: SendMessage(prompt, context)}
	}
}

func RunChatUI() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting chat UI: ", err)
		os.Exit(1)
	}
}
