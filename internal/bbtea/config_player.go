package bbtea

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/e-gloo/orlog/internal/client"
	g "github.com/e-gloo/orlog/internal/client/game"
)

var godNameStyle = lipgloss.NewStyle().Bold(true).Underline(true)

type configPlayerModel struct {
	client            c.Client
	textInput         textinput.Model
	usernameValidated bool
	choices           []g.ClientGod
	cursor            int
	selected          []int
	maxGod            int
	validated         bool
}

func initialConfigPlayerModel(client c.Client) tea.Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return configPlayerModel{
		client:            client,
		textInput:         ti,
		usernameValidated: false,
		choices:           client.GetGameGods(),
		cursor:            0,
		selected:          []int{},
		maxGod:            3,
		validated:         false,
	}
}

func (cp configPlayerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cp configPlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if cp.usernameValidated {
		return cp.handleUpdateChoiceInput(msg)
	} else {
		return cp.handleUpdateTextInput(msg)
	}
}

func (cp configPlayerModel) handleUpdateTextInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if cp.textInput.Value() != "" {
				cp.usernameValidated = true
				// su.client.AddPlayerName(su.textInput.Value())
			}
		default:
			cp.textInput, cmd = cp.textInput.Update(msg)
		}
	}
	return cp, cmd
}

func (cp configPlayerModel) handleUpdateChoiceInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		// The "up" and "k" keys move the cursor up
		case tea.KeyUp:
			if cp.cursor > 0 {
				cp.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case tea.KeyDown:
			if cp.cursor < len(cp.choices)-1 {
				cp.cursor++
			}
		case tea.KeySpace:
			if idx := slices.Index(cp.selected, cp.cursor); idx != -1 {
				cp.selected = slices.Delete(cp.selected, idx, idx+1)
			} else if len(cp.selected) < cp.maxGod {
				cp.selected = append(cp.selected, cp.cursor)
			}
		case tea.KeyEnter:
			if len(cp.selected) == 3 {
				cp.validated = true
				cp.client.AddPlayerName(cp.textInput.Value(), [3]int(cp.selected))
			}
		}
	}
	return cp, cmd
}

func (cp configPlayerModel) View() string {
	if cp.validated {
		return fmt.Sprintf("Get ready %s, game is about to start!", cp.textInput.Value())
	} else if cp.usernameValidated {
		return cp.viewChoiceInput()
	} else {
		return cp.viewTextInput()
	}
}

func (cp configPlayerModel) viewTextInput() string {
	var s string

	if ph.Phase() == c.ConfigPlayer && cp.client.Error() != "" {
		s = cp.client.Error() + "\n"
		s += fmt.Sprintf(
			"Enter your player name:\n\n %s\n",
			cp.textInput.View(),
		)
		return s
	}

	s += fmt.Sprintf(
		"Enter your player name:\n\n %s\n",
		cp.textInput.View(),
	)
	return s
}

func (cp configPlayerModel) viewChoiceInput() string {
	var s string
	s += "Select 3 gods:\n"

	// Iterate over our choices
	for i, choice := range cp.choices {
		cursor := " "
		if cp.cursor == i {
			cursor = ">"
		}

		selection := "[ ]"
		if len(cp.selected) > 0 && slices.Contains(cp.selected, i) {
			selection = "[x]"
		}

		// Render the row
		s += fmt.Sprintf("\t%s %s\t%s: %s\n", cursor, selection, godNameStyle.Render(choice.Name), choice.Description)
	}

	return s
}
