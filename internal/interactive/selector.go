package interactive

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type SelectOption struct {
	ID      string
	Display string
}

type selectorModel struct {
	title       string
	placeholder string
	textInput   textinput.Model
	options     []SelectOption
	filtered    []SelectOption
	cursor      int
	selected    string
	width       int
	height      int
	quitting    bool
	cancelled   bool
}

func InteractiveSelect(title, placeholder string, options []SelectOption, currentValue string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 60

	cursor := 0
	for i, opt := range options {
		if opt.ID == currentValue {
			cursor = i
			break
		}
	}

	m := selectorModel{
		title:       title,
		placeholder: placeholder,
		textInput:   ti,
		options:     options,
		filtered:    options,
		cursor:      cursor,
		width:       80,
		height:      20,
	}

	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	res := final.(selectorModel)
	if res.cancelled {
		return "", fmt.Errorf("selection cancelled")
	}
	return res.selected, nil
}

func (m selectorModel) Init() tea.Cmd { return textinput.Blink }

func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			m.quitting = true
			m.cancelled = true
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor].ID
				m.quitting = true
				return m, tea.Quit
			}
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			m.filtered = filterOptions(m.options, m.textInput.Value())
			if m.cursor >= len(m.filtered) {
				m.cursor = len(m.filtered) - 1
			}
			if m.cursor < 0 {
				m.cursor = 0
			}
			return m, cmd
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m selectorModel) View() string {
	if m.quitting {
		return ""
	}
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")
	countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	b.WriteString(countStyle.Render(fmt.Sprintf("Showing %d of %d options", len(m.filtered), len(m.options))))
	b.WriteString("\n\n")
	sel := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	norm := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	max := 10
	start := m.cursor - max/2
	if start < 0 {
		start = 0
	}
	end := start + max
	if end > len(m.filtered) {
		end = len(m.filtered)
		start = end - max
		if start < 0 {
			start = 0
		}
	}
	for i := start; i < end; i++ {
		opt := m.filtered[i]
		if i == m.cursor {
			b.WriteString(sel.Render("> " + opt.Display))
		} else {
			b.WriteString(norm.Render("  " + opt.Display))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	b.WriteString(help.Render("↑/↓: navigate • Enter: select • Esc: cancel"))
	return b.String()
}

func filterOptions(options []SelectOption, search string) []SelectOption {
	if search == "" {
		return options
	}
	s := strings.ToLower(search)
	var out []SelectOption
	for _, o := range options {
		if strings.Contains(strings.ToLower(o.ID), s) || strings.Contains(strings.ToLower(o.Display), s) {
			out = append(out, o)
		}
	}
	return out
}
