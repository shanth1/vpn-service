package bubbletea

// import (
// 	"fmt"

// 	"github.com/charmbracelet/bubbles/help"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// var (
// 	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
// 	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
// )

// type model struct {
// 	textInput textinput.Model
// 	choices   []string
// 	cursor    int
// 	selected  map[int]struct{}
// 	help      help.Model
// 	keymap    keymap
// }

// type keymap struct{}

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {

// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit

// 		case "up", "k":
// 			if m.cursor > 0 {
// 				m.cursor--
// 			}

// 		case "down", "j":
// 			if m.cursor < len(m.choices)-1 {
// 				m.cursor++
// 			}

// 		case "tab":
// 			m.textInput.Focus()

// 		case "enter", " ":
// 			_, ok := m.selected[m.cursor]
// 			if ok {
// 				delete(m.selected, m.cursor)
// 			} else {
// 				m.selected[m.cursor] = struct{}{}
// 			}
// 		}

// 	}

// 	var cmd tea.Cmd
// 	m.textInput, cmd = m.textInput.Update(msg)
// 	return m, cmd
// }

// func (m model) View() string {
// 	s := "What should we buy at the market?\n\n"

// 	s += fmt.Sprintf("%s\n\n", keywordStyle.Render("KEYWORD"))
// 	s += fmt.Sprintf("%s\n\n", helpStyle.Render("HELP"))

// 	for i, choice := range m.choices {
// 		cursor := " "
// 		if m.cursor == i {
// 			cursor = ">"
// 		}

// 		checked := " "
// 		if _, ok := m.selected[i]; ok {
// 			checked = "x"
// 		}

// 		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
// 	}

// 	s += fmt.Sprintf("\n%s\n", m.textInput.View())

// 	s += "\nPress q to quit.\n"

// 	return s
// }

// func initialModel() model {
// 	ti := textinput.New()
// 	ti.Placeholder = "repository"
// 	ti.Prompt = "charmbracelet/"
// 	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
// 	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
// 	// ti.Focus()
// 	ti.CharLimit = 50
// 	ti.Width = 20
// 	ti.ShowSuggestions = true

// 	h := help.New()

// 	km := keymap{}

// 	return model{
// 		textInput: ti,
// 		choices:   []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
// 		selected:  make(map[int]struct{}),
// 		help:      h,
// 		keymap:    km,
// 	}
// }
