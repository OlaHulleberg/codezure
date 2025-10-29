package interactive

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InteractiveInput shows a prompt and returns the entered text.
func InteractiveInput(title, placeholder, initial string) (string, error) {
	return runTextInput(title, placeholder, initial, false)
}

// InteractivePassword shows a prompt and returns the entered secret (masked input).
func InteractivePassword(title, placeholder string) (string, error) {
	return runTextInput(title, placeholder, "", true)
}

type inputModel struct {
	title     string
	ti        textinput.Model
	quitting  bool
	cancelled bool
}

func runTextInput(title, placeholder, initial string, secret bool) (string, error) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(initial)
	ti.Focus()
	ti.Width = 60
	if secret {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	m := inputModel{title: title, ti: ti}
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	res := final.(inputModel)
	if res.cancelled {
		return "", fmt.Errorf("input cancelled")
	}
	return res.ti.Value(), nil
}

func (m inputModel) Init() tea.Cmd { return textinput.Blink }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			m.quitting = true
			m.cancelled = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	if m.quitting {
		return ""
	}
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	return titleStyle.Render(m.title) + "\n" + m.ti.View() + "\n\n" + helpStyle.Render("Enter: confirm • Esc: cancel")
}
