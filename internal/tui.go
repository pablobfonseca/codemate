package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA500")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF8700")).
			Padding(0, 1)

	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			BorderLeft(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			PaddingLeft(1)

	aiMsgStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			PaddingLeft(1)

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFA500")).
			Padding(0, 1)

	timestampStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Italic(true).
			MarginLeft(1).
			MarginBottom(1).
			MarginTop(1)
)

type Message struct {
	content   string
	isUser    bool
	timestamp time.Time
}

type model struct {
	messages  []Message
	textInput textinput.Model
	viewport  viewport.Model
	spinner   spinner.Model
	loading   bool
	width     int
	height    int
	ready     bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Ask something about your code..."
	ti.Focus()
	ti.Width = 80
	ti.CharLimit = 500

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))

	return model{
		textInput: ti,
		messages:  []Message{},
		spinner:   s,
		loading:   false,
		ready:     false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
		tea.EnterAltScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.loading {
				return m, nil
			}

			if strings.TrimSpace(m.textInput.Value()) == "" {
				return m, nil
			}

			userMessage := m.textInput.Value()
			m.messages = append(m.messages, Message{
				content:   userMessage,
				isUser:    true,
				timestamp: time.Now(),
			})
			m.textInput.Reset()
			m.loading = true
			m.updateViewport()

			return m, tea.Batch(
				m.spinner.Tick,
				requestChat(userMessage),
			)
		case "pgup":
			m.viewport.HalfViewUp()
		case "pgdown":
			m.viewport.HalfViewDown()
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		headerHeight := 3
		footerHeight := 6

		m.viewport = viewport.New(msg.Width-4, msg.Height-headerHeight-footerHeight)
		m.viewport.SetContent(m.contentAsString())
		m.viewport.YPosition = headerHeight
		m.textInput.Width = msg.Width - 10
		m.ready = true

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case streamChunkMsg:
		if len(m.messages) == 0 || m.messages[len(m.messages)-1].isUser {
			m.messages = append(m.messages, Message{
				content:   msg.chunk,
				isUser:    false,
				timestamp: time.Now(),
			})
		} else {
			lastIndex := len(m.messages) - 1
			m.messages[lastIndex].content += msg.chunk
		}
		m.updateViewport()
		return m, nil

	case streamCompleteMsg:
		m.loading = false
		if len(m.messages) >= 2 {
			userMsg := m.messages[len(m.messages)-2].content
			aiMsg := m.messages[len(m.messages)-1].content
			SaveMessage(userMsg, aiMsg)
		}
		return m, nil
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)

	cmds = append(cmds, tiCmd, vpCmd, spCmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var sb strings.Builder

	title := titleStyle.Render("Welcome to Codemate! - Your Code Assistant")
	sb.WriteString(title + "\n\n")

	contentView := m.viewport.View()
	sb.WriteString(contentView)
	sb.WriteString("\n\n")

	inputPrompt := ""
	if m.loading {
		inputPrompt = m.spinner.View() + "Thinking...\n\n"
	}

	inputBar := inputPrompt + inputStyle.Render(m.textInput.View())
	sb.WriteString(inputBar)

	return appStyle.Render(sb.String())
}

type streamChunkMsg struct {
	chunk string
}

type streamCompleteMsg struct{}

var chunkChannel = make(chan string, 100)
var doneChannel = make(chan bool, 1)

func requestChat(prompt string) tea.Cmd {
	return func() tea.Msg {
		context, _ := GetProjectContext()

		SetStreamCallback(func(chunk string, done bool) {
			if done {
				doneChannel <- true
			} else {
				chunkChannel <- chunk
			}
		})

		go func() {
			SendMessage(prompt, context)
		}()

		go func() {
			for {
				select {
				case chunk := <-chunkChannel:
					program.Send(streamChunkMsg{chunk: chunk})
				case <-doneChannel:
					program.Send(streamCompleteMsg{})
					return
				}
			}
		}()

		return nil
	}
}

func formatMessage(msg Message) string {
	if msg.isUser {
		header := userMsgStyle.Render("You")
		content := userMsgStyle.Render(msg.content)
		return header + "\n" + content
	} else {
		rendered, _ := glamour.Render(msg.content, "dark")
		header := aiMsgStyle.Render("Codemate")
		content := aiMsgStyle.Render(rendered)
		return header + "\n" + content
	}
}

func (m *model) contentAsString() string {
	var sb strings.Builder

	for _, msg := range m.messages {
		sb.WriteString(formatMessage(msg))
		sb.WriteString("\n\n")
	}

	return sb.String()
}

func (m *model) updateViewport() {
	m.viewport.SetContent(m.contentAsString())
	m.viewport.GotoBottom()
}

var program *tea.Program

func RunChatUI() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	program = p

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
