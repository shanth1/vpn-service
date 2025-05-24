package users

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tuicommon "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/common"
	"github.com/shanth1/vpn-service/internal/core/domain"
	"github.com/shanth1/vpn-service/internal/core/ports"
)

const TabTitle = "Users"

type viewState int

const (
	stateListLoading viewState = iota
	stateListView
	stateDetailsLoading
	stateDetailsView
	stateCreateForm
	stateSubmittingForm
)

type Model struct {
	core   ports.PrimaryPort
	width  int
	height int
	state  viewState

	spinner    spinner.Model
	err        error
	successMsg string

	users             []userListItem
	selectedUserIndex int

	currentUserDetails *domain.User
	currentUserTraffic *domain.TrafficInfo

	formInputs     []textinput.Model
	formFocusIndex int
}

func New(core ports.PrimaryPort) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(tuicommon.AccentColor)

	// TODO: change to map
	placeholders := []string{"Telegram Username", "Name", "Surname", "Email"}
	inputs := make([]textinput.Model, len(placeholders))
	charLimit := 32
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].Width = charLimit
		inputs[i].Placeholder = placeholders[i]
		inputs[i].CharLimit = charLimit
		if i == 0 {
			inputs[i].Focus()
		}
		inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(tuicommon.AccentColor)
		inputs[i].TextStyle = lipgloss.NewStyle()
	}

	return Model{
		core:       core,
		spinner:    s,
		state:      stateListLoading,
		formInputs: inputs,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchAllUsersCmd(m.core))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if _, ok := msg.(tea.KeyMsg); ok {
		m.err = nil
		m.successMsg = ""
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.state {
		case stateListView:
			cmd = m.handleListViewKeys(msg)
		case stateDetailsView:
			cmd = m.handleDetailsViewKeys(msg)
		case stateCreateForm:
			var formCmd tea.Cmd
			m, formCmd = m.handleCreateFormKeys(msg)
			cmds = append(cmds, formCmd)
		}

	// Async commands
	case allUsersLoadedMsg:
		m.spinner.Tick()
		if msg.err != nil {
			m.err = fmt.Errorf("failed to load users: %w", msg.err)
			m.state = stateListView
		} else {
			m.users = []userListItem{}
			for _, u := range msg.users {
				m.users = append(m.users, userListItem{user: u})
			}
			m.selectedUserIndex = 0
			if len(m.users) == 0 {
				m.selectedUserIndex = -1
			}
			m.state = stateListView
		}
	case userDetailsLoadedMsg:
		m.spinner.Tick()
		if msg.err != nil {
			m.err = fmt.Errorf("failed to load user details: %w", msg.err)
			m.state = stateDetailsView
		} else {
			m.currentUserTraffic = msg.traffic
			m.state = stateDetailsView
		}
	case userDeletedMsg:
		m.spinner.Tick()
		if msg.err != nil {
			m.err = fmt.Errorf("failed to delete user %s: %w", msg.userID, msg.err)
			m.state = stateDetailsView
		} else {
			m.successMsg = fmt.Sprintf("User %s deleted successfully.", msg.userID)
			m.currentUserDetails = nil
			m.currentUserTraffic = nil
			m.state = stateListLoading
			cmd = fetchAllUsersCmd(m.core)
		}
	case userCreatedMsg:
		m.spinner.Tick()
		if msg.err != nil {
			m.err = fmt.Errorf("failed to create user: %w", msg.err)
			m.state = stateCreateForm
		} else {
			m.successMsg = fmt.Sprintf("User %s created successfully.", msg.user.TG)
			// Clear form
			for i := range m.formInputs {
				m.formInputs[i].SetValue("")
			}
			m.state = stateListLoading
			cmd = fetchAllUsersCmd(m.core)
		}

	case spinner.TickMsg:
		if m.state == stateListLoading || m.state == stateDetailsLoading || m.state == stateSubmittingForm {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			cmds = append(cmds, spinCmd)
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing users tab..."
	}

	var viewContent string

	switch m.state {
	case stateListLoading:
		loadingText := fmt.Sprintf("%s Loading users...", m.spinner.View())
		viewContent = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, loadingText)

	case stateListView:
		viewContent = m.renderListView()

	case stateDetailsLoading:
		loadingText := fmt.Sprintf("%s Loading user details...", m.spinner.View())
		leftPane := m.renderUserListPane(m.width/2-1, m.height)
		rightPane := lipgloss.Place(m.width/2-1, m.height, lipgloss.Center, lipgloss.Center, loadingText)
		viewContent = lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	case stateDetailsView:
		viewContent = m.renderDetailsView()

	case stateCreateForm:
		viewContent = m.renderCreateFormView()

	case stateSubmittingForm:
		submittingText := fmt.Sprintf("%s Submitting...", m.spinner.View())
		viewContent = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, submittingText)
	}

	var finalMsg string
	if m.err != nil {
		finalMsg = tuicommon.ErrorStyle.Render("Error: " + m.err.Error())
	} else if m.successMsg != "" {
		finalMsg = tuicommon.SuccessStyle.Render(m.successMsg)
	}

	if finalMsg != "" {
		return lipgloss.JoinVertical(lipgloss.Left, viewContent, finalMsg)
	}
	return viewContent
}
